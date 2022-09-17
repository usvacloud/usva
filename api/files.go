package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romeq/usva/config"
	"github.com/romeq/usva/dbengine"
	"golang.org/x/crypto/bcrypt"
)

var errAuthMissing = errors.New("missing authentication in protected file")

func uploadFile(ctx *gin.Context, uploadOptions *config.Files) {
	// retrieve file from request
	f, err := ctx.FormFile("file")
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	// save file
	filename := uuid.New().String() + path.Ext(f.Filename)

	if int(f.Size)/1000000 > uploadOptions.MaxSize {
		ctx.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "File is too big",
		})
		return
	}

	// decode key as base64
	var hash []byte
	pwd := strings.TrimSpace(ctx.PostForm("password"))
	if len(pwd) > 0 {
		decodedkey, err := base64.RawStdEncoding.DecodeString(pwd)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Failed to decode key as base64",
			})
			return
		}
		decodedkey = []byte(strings.TrimSpace(string(decodedkey)))

		// generate the password hash
		if len(decodedkey) > 0 {
			hash, err = bcrypt.GenerateFromPassword(decodedkey, 12)
			if err != nil {
				setErrResponse(ctx, err)
				return
			}
		}
	}

	// Append file metadata into database
	err = dbengine.InsertFile(dbengine.File{
		Filename:    filename,
		FileSize:    int(f.Size),
		Password:    string(hash),
		IsEncrypted: len(pwd) > 0 && ctx.PostForm("didClientEncrypt") == "yes",
		UploadDate:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Println(err)
		setErrResponse(ctx, err)
		return
	}

	err = ctx.SaveUploadedFile(f, path.Join(uploadOptions.UploadsDir, filename))
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded",
		"filename": filename,
	})
}

func fileInformation(ctx *gin.Context) {
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

	f, err := dbengine.GetFile(filename)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	pwd, err := dbengine.GetPasswordHash(filename)
	if err != nil {
		setErrResponse(ctx, err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"filename":   f.Filename,
		"size":       f.FileSize,
		"uploadDate": f.UploadDate,
		"viewCount":  f.ViewCount,
		"locked":     pwd != "",
		"encrypted":  f.IsEncrypted,
	})
}

func downloadFile(ctx *gin.Context, downloadOptions *config.Files) {
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

	err := dbengine.IncrementFileViewCount(filename)
	if err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.FileAttachment(path.Join(downloadOptions.UploadsDir, filename), filename)
}

func deleteFile(ctx *gin.Context, fileOptions *config.Files) {
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

	if err := dbengine.DeleteFile(filename); err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "file deleted!",
	})
}

func authorizeRequest(ctx *gin.Context, filename string) (success bool) {
	pwdhash, err := dbengine.GetPasswordHash(filename)
	if err != nil {
		setErrResponse(ctx, err)
		return false
	}

	if len(pwdhash) > 0 {
		pwd, err := parseHeaderPassword(ctx)
		if err != nil {
			setErrResponse(ctx, err)
			return false
		}

		trimmedpwd := strings.TrimSpace(pwd)
		err = bcrypt.CompareHashAndPassword([]byte(pwdhash), []byte(trimmedpwd))
		if err != nil {
			setErrResponse(ctx, err)
			return false
		}
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

func setErrResponse(ctx *gin.Context, err error) {
	errorMessage, status := "request failed", http.StatusBadRequest

	switch err {
	case sql.ErrNoRows:
		errorMessage, status = "file not found", http.StatusNotFound
	case bcrypt.ErrMismatchedHashAndPassword:
		errorMessage, status = "password is invalid", http.StatusForbidden
	case errAuthMissing:
		errorMessage, status = err.Error(), http.StatusUnauthorized
	default:
		log.Println("error: ", err.Error())
	}

	ctx.AbortWithStatusJSON(status, gin.H{
		"error": errorMessage,
	})
}
