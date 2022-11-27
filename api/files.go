package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romeq/usva/api/middleware"
	"github.com/romeq/usva/db"
	"github.com/romeq/usva/dbengine"
	"github.com/romeq/usva/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	errAuthMissing = errors.New("missing authentication in a protected file")
	errAuthFailed  = errors.New("authorization succeeded but cookies were not set")
	errInvalidBody = errors.New("invalid request body")
)

func UploadFile(lmt *middleware.Ratelimiter, uploadOptions *APIConfiguration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// retrieve file from request
		f, err := ctx.FormFile("file")
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

		// generate name for the uploaded file
		filename := uuid.NewString() + path.Ext(f.Filename)

		if uploadOptions.MaxSingleUploadSize > 0 && uint64(f.Size) > uploadOptions.MaxSingleUploadSize {
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
		err = dbengine.DB.NewFile(ctx, db.NewFileParams{
			FileUuid:    filename,
			Title:       sql.NullString{String: title, Valid: title != ""},
			Passwdhash:  sql.NullString{String: string(hash), Valid: string(hash) != ""},
			Uploader:    sql.NullString{String: apiid, Valid: apiid != ""},
			AccessToken: uuid.NewString(),
			Isencrypted: len(pwd) > 0 && ctx.PostForm("didClientEncrypt") == "yes",
			UploadDate:  time.Now().Format(time.RFC3339),
		})
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

		abspath, err := filepath.Abs(uploadOptions.UploadsDir)
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

		err = ctx.SaveUploadedFile(f, path.Join(abspath, filename))
		if err != nil {
			_ = dbengine.DB.DeleteFile(ctx, filename)
			setErrResponse(ctx, err)
			return
		}

		lmt.AppendClientUploads(apiid, middleware.ClientUpload{
			Size: f.Size,
			Time: time.Now(),
		})
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "file uploaded",
			"filename": filename,
		})
	}
}

// => /file/info?filename=<uuid>
func FileInformation(uploadInformation *APIConfiguration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filename, filenameGiven := ctx.GetQuery("filename")
		if !filenameGiven {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Filename not given",
			})
			return
		}

		if !authorizeRequest(ctx, filename) {
			return
		}

		f, err := dbengine.DB.FileInformation(ctx, filename)
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

		pwd, err := dbengine.DB.GetPasswordHash(ctx, filename)
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

		filesize, err := utils.FileSize(path.Join(uploadInformation.UploadsDir, filename))
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
}

func DownloadFile(uploadInformation *APIConfiguration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filename, filenameGiven := ctx.GetQuery("filename")
		if !filenameGiven {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "filename not given",
			})
			return
		}

		if !authorizeRequest(ctx, filename) {
			return
		}

		err := dbengine.DB.UpdateViewCount(ctx, filename)
		if err != nil {
			setErrResponse(ctx, err)
			return
		}

		ctx.FileAttachment(path.Join(uploadInformation.UploadsDir, filename), filename)
	}
}

func ReportFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		err = dbengine.DB.NewReport(ctx, db.NewReportParams{
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
}

func DeleteFile(fileOptions *APIConfiguration) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		filename, filenameGiven := ctx.GetQuery("filename")
		if !filenameGiven {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "filename not given",
			})
			return
		}

		if !authorizeRequest(ctx, filename) {
			return
		}

		if err := os.Remove(path.Join(fileOptions.UploadsDir, filename)); err != nil {
			setErrResponse(ctx, err)
			return
		}

		if err := dbengine.DB.DeleteFile(ctx, filename); err != nil {
			setErrResponse(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "file deleted!",
		})
	}
}

// Functions to help with most common tasks
func authorizeRequest(ctx *gin.Context, filename string) (success bool) {
	pwdhash, err := dbengine.DB.GetPasswordHash(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return false
	}

	if pwdhash.Valid {
		fileauthcookie := fmt.Sprintf("usva-auth-%s", filename)
		authcookieValue, _ := ctx.Cookie(fileauthcookie)

		at, err := dbengine.DB.GetAccessToken(ctx, filename)
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

		ctx.SetCookie(fileauthcookie, at, 3600, "/", "localhost", false, false)
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

	switch err {
	case sql.ErrNoRows:
		errorMessage, status = "file not found", http.StatusNotFound
	case bcrypt.ErrMismatchedHashAndPassword:
		errorMessage, status = "password is invalid", http.StatusForbidden
	case errInvalidBody:
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errAuthMissing:
		errorMessage, status = err.Error(), http.StatusUnauthorized
	case errEmptyResponse:
		errorMessage, status = err.Error(), http.StatusNoContent
	default:
		log.Println("error: ", err.Error())
	}

	ctx.AbortWithStatusJSON(status, gin.H{
		"error": errorMessage,
	})
}
