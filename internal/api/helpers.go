package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/internal/cryptography"
	"golang.org/x/crypto/bcrypt"
)

var (
	errNotFound      = errors.New("resource not found")
	errEmptyResponse = errors.New("empty response")
	errTooBigBody    = errors.New("request body too big")
)

func prettybytes(bytes uint64) gin.H {
	return gin.H{
		"bytes":     bytes,
		"kilobytes": bytes / 1000,
		"megabytes": bytes / 1000 / 1000,
		"gigabytes": bytes / 1000 / 1000 / 1000,
	}
}

func (s *Server) RestrictionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSingleUploadSize":  prettybytes(s.api.MaxSingleUploadSize),
		"maxDailyUploadSize":   prettybytes(s.api.MaxUploadSizePerDay),
		"maxEncryptedFileSize": prettybytes(s.api.MaxEncryptableFileSize),
		"filePersistDuration": gin.H{
			"seconds": s.api.FilePersistDuration.Seconds(),
			"hours":   s.api.FilePersistDuration.Hours(),
			"days":    s.api.FilePersistDuration.Hours() / 24,
		},
	})
}

func (s *Server) NotFoundHandler(ctx *gin.Context) {
	setErrResponse(ctx, errNotFound)
}

// setErrResponse helper for providing standard error messages in return
func setErrResponse(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	errorMessage, status := "request failed", http.StatusBadRequest

	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		errorMessage, status = "password is invalid", http.StatusForbidden
	case errors.Is(err, cryptography.ErrPasswordTooShort):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, errTooBigBody):
		errorMessage, status = err.Error(), http.StatusRequestEntityTooLarge
	case errors.Is(err, cryptography.ErrPasswordTooLong):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, errInvalidBody):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, errAuthMissing):
		errorMessage, status = err.Error(), http.StatusUnauthorized
	case errors.Is(err, errNotFound):
		errorMessage, status = err.Error(), http.StatusNotFound
	case errors.Is(err, sql.ErrNoRows):
		errorMessage, status = err.Error(), http.StatusNoContent
	case errors.Is(err, errEmptyResponse):
		errorMessage, status = err.Error(), http.StatusNoContent

	default:
		log.Println("error: ", err.Error())
	}

	ctx.AbortWithStatusJSON(status, gin.H{
		"error": errorMessage,
	})
}
