package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romeq/usva/db"
	"github.com/romeq/usva/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	errAuthMissing = errors.New("missing authentication in a protected file")
	errAuthFailed  = errors.New("authorization succeeded but cookies were not set")
	errInvalidBody = errors.New("invalid request body")
)

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

	apiid := ctx.Writer.Header().Get("Api-Identifier")
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

	ctx.String(http.StatusOK, fmt.Sprintf("%s://%s/file?filename=%s",
		protocol, CookieDomain, filename))
}

func (s *Server) UploadFile(ctx *gin.Context) {
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

	var hash []byte
	pwd := strings.TrimSpace(ctx.PostForm("password"))
	if len(pwd) > 0 {
		if len(pwd) > 512 {
			setErrResponse(ctx, errInvalidBody)
			return
		}

		decodedkey, err := base64.RawStdEncoding.DecodeString(pwd)
		if err != nil {
			setErrResponse(ctx, errInvalidBody)
			return
		}

		decodedkey = []byte(strings.TrimSpace(string(decodedkey)))
		if len(decodedkey) > 0 {
			hash, err = bcrypt.GenerateFromPassword(decodedkey, 12)
			if err != nil {
				setErrResponse(ctx, err)
				return
			}
		}
	}

	// Append file metadata into database
	apiid := ctx.Writer.Header().Get("Api-Identifier")
	title := ctx.Request.FormValue("title")
	err = s.db.NewFile(ctx, db.NewFileParams{
		FileUuid:    filename,
		Title:       sql.NullString{String: title, Valid: title != ""},
		Passwdhash:  sql.NullString{String: string(hash), Valid: string(hash) != ""},
		Uploader:    sql.NullString{String: apiid, Valid: apiid != ""},
		AccessToken: uuid.NewString(),
		Isencrypted: len(pwd) > 0 && ctx.PostForm("didClientEncrypt") == "yes",
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
		"encrypted":  f.Isencrypted,
	})
}

func (s *Server) DownloadFile(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "filename not given",
		})
		return
	}

	if !s.authorizeRequest(ctx, filename) {
		return
	}

	if err := s.db.UpdateLastSeen(ctx, filename); err != nil {
		setErrResponse(ctx, err)
	}

	err := s.db.UpdateViewCount(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.FileAttachment(path.Join(s.api.UploadsDir, filename), filename)
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

// Functions to help with most common tasks
func (s *Server) authorizeRequest(ctx *gin.Context, filename string) bool {
	pwdhash, err := s.db.GetPasswordHash(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return false
	}

	if pwdhash.Valid {
		fileauthcookie := fmt.Sprintf("usva-auth-%s", filename)
		authcookieValue, _ := ctx.Cookie(fileauthcookie)

		at, err := s.db.GetAccessToken(ctx, filename)
		if err != nil {
			setErrResponse(ctx, errAuthFailed)
			return false
		}

		if authcookieValue == at {
			return true
		}

		pwd, err := parseHeaderPassword(ctx)
		if err != nil {
			setErrResponse(ctx, err)
			return false
		}

		trimmedpwd := strings.TrimSpace(pwd)
		err = bcrypt.CompareHashAndPassword([]byte(pwdhash.String), []byte(trimmedpwd))
		if err != nil {
			setErrResponse(ctx, err)
			return false
		}

		ctx.SetCookie(fileauthcookie, at, s.api.CookieSaveTime, "/",
			CookieDomain, s.api.UseSecureCookie, true,
		)
	}

	return true
}

func parseHeaderPassword(ctx *gin.Context) (string, error) {
	authheader := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
	if len(authheader) < 2 {
		return "", errAuthMissing
	}

	p, err := base64.RawStdEncoding.DecodeString(authheader[1])
	if err != nil {
		return "", err
	}

	return string(p), nil
}

// helper for providing standard error messages in return
func setErrResponse(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	errorMessage, status := "request failed", http.StatusBadRequest

	switch {
	case errors.Is(err, sql.ErrNoRows):
		errorMessage, status = "file not found", http.StatusNotFound
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		errorMessage, status = "password is invalid", http.StatusForbidden
	case errors.Is(err, errInvalidBody):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, errAuthMissing):
		errorMessage, status = err.Error(), http.StatusUnauthorized
	case errors.Is(err, errEmptyResponse):
		errorMessage, status = err.Error(), http.StatusNoContent
	default:
		log.Println("error: ", err.Error())
	}

	ctx.AbortWithStatusJSON(status, gin.H{
		"error": errorMessage,
	})
}
