package file

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
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/cmd/webserver/handlers/auth"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/internal/utils"
	"github.com/romeq/usva/pkg/cryptography"
	"github.com/romeq/usva/pkg/ratelimit"
)

type Handler struct {
	db                *db.Queries
	api               *handlers.Configuration
	encryptionKeySize uint32
	auth              *auth.Handler
}

func NewFileHandler(s *handlers.Server, authHandler *auth.Handler) *Handler {
	return &Handler{
		db:                s.DB,
		api:               s.Config,
		encryptionKeySize: s.EncKeySize,
		auth:              authHandler,
	}
}

// UploadFileSimple is a simple wrapper around ctx.SaveUploadedFile to support
// an upload with a very very simple curl request (curl -Ld = )
func (s *Handler) UploadFileSimple(ctx *gin.Context) {
	f, err := ctx.FormFile("file")
	if err != nil {
		handlers.SetErrResponse(ctx, err)
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
		handlers.SetErrResponse(ctx, err)
		return
	}

	abspath, err := filepath.Abs(s.api.UploadsDir)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	err = ctx.SaveUploadedFile(f, path.Join(abspath, filename))
	if err != nil {
		_ = s.db.DeleteFile(ctx, filename)
		handlers.SetErrResponse(ctx, err)
		return
	}

	protocol := "http"
	if s.api.UseSecureCookie {
		protocol = "https"
	}

	ctx.String(http.StatusOK, fmt.Sprintf("%s://%s/file/?filename=%s", protocol, s.api.APIDomain, filename))
}

// UploadFile
func (s *Handler) UploadFile(ctx *gin.Context) {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	filename := uuid.NewString() + path.Ext(formFile.Filename)

	if s.api.MaxSingleUploadSize > 0 && uint64(formFile.Size) > s.api.MaxSingleUploadSize {
		handlers.SetErrResponse(ctx, handlers.ErrTooBigBody)
		return
	}

	formpwd := strings.TrimSpace(ctx.PostForm("password"))
	pwd, err := base64.RawStdEncoding.DecodeString(formpwd)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}
	password := strings.TrimSpace(string(pwd))

	formFileHandle, err := formFile.Open()
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	absUploads, err := filepath.Abs(s.api.UploadsDir)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	file, err := os.OpenFile(path.Join(absUploads, filename), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
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
			handlers.SetErrResponse(ctx, err)
			return
		}

		iv, err = cryptography.EncryptStream(file, formFileHandle, encryptionKey)
		if err != nil {
			handlers.SetErrResponse(ctx, err)
			return
		}

	case !requirementsMet && confirmation:
		handlers.SetErrResponse(ctx, handlers.ErrInvalidBody)
		return

	default:
		if _, err := io.Copy(file, formFileHandle); err != nil {
			handlers.SetErrResponse(ctx, err)
			return
		}
	}

	hash, err := auth.BCryptPasswordHash([]byte(password))
	if err != nil {
		handlers.SetErrResponse(ctx, err)
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
		handlers.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded",
		"filename": filename,
	})
}

// => /file/info?filename=<uuid>
func (s *Handler) FileInformation(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Filename not given",
		})
		return
	}

	if !s.auth.AuthorizeRequest(ctx, filename) {
		return
	}

	f, err := s.db.FileInformation(ctx, filename)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	if err = s.db.UpdateLastSeen(ctx, filename); err != nil {
		handlers.SetErrResponse(ctx, err)
	}

	pwd, err := s.db.GetPasswordHash(ctx, filename)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	filesize, err := utils.FileSize(path.Join(s.api.UploadsDir, filename))
	if err != nil {
		handlers.SetErrResponse(ctx, err)
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

func (s *Handler) DownloadFile(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		handlers.SetErrResponse(ctx, errors.New("filename not given"))
		return
	}

	// authorize request
	if !s.auth.AuthorizeRequest(ctx, filename) {
		return
	}

	filepath := path.Join(s.api.UploadsDir, filename)
	fileHandle, err := os.Open(filepath)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	headerPassword, err := s.auth.ParseFilePassword(ctx, filename)
	if err != nil && !errors.Is(err, handlers.ErrAuthMissing) {
		handlers.SetErrResponse(ctx, err)
		return
	}

	encryptionIv, err := s.db.GetDownload(ctx, filename)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	if len(encryptionIv) == 0 {
		ctx.FileAttachment(filepath, path.Base(filepath))
		return
	} else if errors.Is(err, handlers.ErrAuthMissing) {
		handlers.SetErrResponse(ctx, handlers.ErrAuthMissing)
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", `attachment;`)

	ctx.Status(http.StatusOK)

	derivedKey, err := cryptography.DeriveBasicKey([]byte(headerPassword), s.encryptionKeySize)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	err = cryptography.DecryptStream(ctx.Writer, fileHandle, derivedKey, encryptionIv)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}
}

func (s *Handler) ReportFile(ctx *gin.Context) {
	var requestBody struct {
		Filename string `json:"filename"`
		Reason   string `json:"reason"`
	}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	if len(requestBody.Filename) < 36 ||
		!utils.IsBetween(len(requestBody.Reason), 20, 1024) {

		handlers.SetErrResponse(ctx, handlers.ErrInvalidBody)
		return
	}

	err = s.db.NewReport(ctx, db.NewReportParams{
		FileUuid: requestBody.Filename,
		Reason:   requestBody.Reason,
	})
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "thank you! your report has been sent.",
	})
}
