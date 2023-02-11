package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/internal/utils"
	"github.com/romeq/usva/pkg/cryptography"
	"github.com/romeq/usva/pkg/ratelimit"
)

var (
	errAuthMissing = errors.New("missing authentication in a protected file")
	errAuthFailed  = errors.New("authorization succeeded but cookies were not set")
	errInvalidBody = errors.New("invalid request body")
)

// UploadFileSimple is a simple wrapper around ctx.SaveUploadedFile to support
// an upload with a very very simple curl request (curl -Ld = )
func (s *Server) UploadFileSimple(ctx *gin.Context) {
	// retrieve file from request
	f, err := ctx.FormFile("file")
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	// generate name for the uploaded file
	filename := uuid.NewString() + path.Ext(f.Filename)

	if s.api.MaxSingleUploadSize > 0 && uint64(f.Size) > s.api.MaxSingleUploadSize {
		ctx.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "File is too big",
		})
		return
	}

	apiid := ctx.Writer.Header().Get(ratelimit.Headers.Identifier)
	title := ctx.Request.FormValue("title")
	err = s.db.NewFile(ctx, db.NewFileParams{
		FileUuid:    filename,
		Title:       sql.NullString{String: title, Valid: title != ""},
		Uploader:    sql.NullString{String: apiid, Valid: apiid != ""},
		AccessToken: uuid.NewString(),
	})
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	abspath, err := filepath.Abs(s.api.UploadsDir)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	err = ctx.SaveUploadedFile(f, path.Join(abspath, filename))
	if err != nil {
		_ = s.db.DeleteFile(ctx, filename)
		setErrResponse(ctx, err)
		return
	}

	protocol := "http"
	if s.api.UseSecureCookie {
		protocol = "https"
	}

	ctx.String(http.StatusOK, fmt.Sprintf("%s://%s/file/?filename=%s", protocol, s.api.APIDomain, filename))
}

// UploadFile
func (s *Server) UploadFile(ctx *gin.Context) {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	filename := uuid.NewString() + path.Ext(formFile.Filename)

	if s.api.MaxSingleUploadSize > 0 && uint64(formFile.Size) > s.api.MaxSingleUploadSize {
		setErrResponse(ctx, errTooBigBody)
		return
	}

	formpwd := strings.TrimSpace(ctx.PostForm("password"))
	pwd, err := base64.RawStdEncoding.DecodeString(formpwd)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}
	password := strings.TrimSpace(string(pwd))

	formFileHandle, err := formFile.Open()
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	absUploads, err := filepath.Abs(s.api.UploadsDir)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	file, err := os.OpenFile(path.Join(absUploads, filename), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	var (
		iv              = []byte{}
		requirementsMet = len(password) >= 6 &&
			len(password) < 128 &&
			formFile.Size < int64(s.api.MaxEncryptableFileSize)
		confirmation = ctx.PostForm("can_encrypt") == "yes"
	)

	switch {
	case requirementsMet && confirmation:
		encryptionKey, err := cryptography.DeriveBasicKey([]byte(password), s.encryptionKeySize)
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

		iv, err = cryptography.EncryptStream(file, formFileHandle, encryptionKey)
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

	case !requirementsMet && confirmation:
		setErrResponse(ctx, errInvalidBody)
		return

	default:
		if _, err := io.Copy(file, formFileHandle); err != nil {
			setErrResponse(ctx, err)
			return
		}
	}

	hash, err := bcryptPasswordHash([]byte(password))
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	apiid := ctx.Writer.Header().Get(ratelimit.Headers.Identifier)
	title := ctx.Request.FormValue("title")
	err = s.db.NewFile(ctx, db.NewFileParams{
		FileUuid:     filename,
		Title:        sql.NullString{String: title, Valid: title != ""},
		Passwdhash:   sql.NullString{String: string(hash), Valid: string(hash) != ""},
		Uploader:     sql.NullString{String: apiid, Valid: apiid != ""},
		EncryptionIv: iv,
		AccessToken:  uuid.NewString(),
	})
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded",
		"filename": filename,
	})
}

// => /file/info?filename=<uuid>
func (s *Server) FileInformation(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Filename not given",
		})
		return
	}

	if !s.authorizeRequest(ctx, filename) {
		return
	}

	f, err := s.db.FileInformation(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	if err = s.db.UpdateLastSeen(ctx, filename); err != nil {
		setErrResponse(ctx, err)
	}

	pwd, err := s.db.GetPasswordHash(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	filesize, err := utils.FileSize(path.Join(s.api.UploadsDir, filename))
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"filename":   f.FileUuid,
		"size":       filesize,
		"title":      f.Title,
		"uploadDate": f.UploadDate,
		"viewCount":  f.Viewcount,
		"locked":     pwd.Valid,
		"encrypted":  f.Encrypted,
	})
}

func (s *Server) DownloadFile(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		setErrResponse(ctx, errors.New("filename not given"))
		return
	}

	// authorize request
	if !s.authorizeRequest(ctx, filename) {
		return
	}

	filepath := path.Join(s.api.UploadsDir, filename)
	fileHandle, err := os.Open(filepath)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	headerPassword, err := s.parseFilePassword(ctx, filename)
	if err != nil && !errors.Is(err, errAuthMissing) {
		setErrResponse(ctx, err)
		return
	}

	encryptionIv, err := s.db.GetDownload(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	if len(encryptionIv) == 0 {
		ctx.FileAttachment(filepath, path.Base(filepath))
		return
	} else if errors.Is(err, errAuthMissing) {
		setErrResponse(ctx, errAuthMissing)
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", `attachment;`)

	ctx.Status(http.StatusOK)

	derivedKey, err := cryptography.DeriveBasicKey([]byte(headerPassword), s.encryptionKeySize)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	err = cryptography.DecryptStream(ctx.Writer, fileHandle, derivedKey, encryptionIv)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}
}

func (s *Server) ReportFile(ctx *gin.Context) {
	var requestBody struct {
		Filename string `json:"filename"`
		Reason   string `json:"reason"`
	}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	if len(requestBody.Filename) < 36 ||
		!utils.IsBetween(len(requestBody.Reason), 20, 1024) {

		setErrResponse(ctx, errInvalidBody)
		return
	}

	err = s.db.NewReport(ctx, db.NewReportParams{
		FileUuid: requestBody.Filename,
		Reason:   requestBody.Reason,
	})
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "thank you! your report has been sent.",
	})
}
