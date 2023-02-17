package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/pkg/cryptography"
	"golang.org/x/crypto/bcrypt"
)

func jsonbytes(bytes uint64) gin.H {
	return gin.H{
		"bytes":     bytes,
		"kilobytes": bytes / 1000,
		"megabytes": bytes / 1000 / 1000,
		"gigabytes": bytes / 1000 / 1000 / 1000,
	}
}

func (s *Server) RestrictionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSingleUploadSize":  jsonbytes(s.Config.MaxSingleUploadSize),
		"maxDailyUploadSize":   jsonbytes(s.Config.MaxUploadSizePerDay),
		"maxEncryptedFileSize": jsonbytes(s.Config.MaxEncryptableFileSize),
		"filePersistDuration": gin.H{
			"seconds": s.Config.FilePersistDuration.Seconds(),
			"hours":   s.Config.FilePersistDuration.Hours(),
			"days":    s.Config.FilePersistDuration.Hours() / 24,
		},
	})
}

func (s *Server) NotFoundHandler(ctx *gin.Context) {
	SetErrResponse(ctx, ErrNotFound)
}

// SetErrResponse helper for providing standard error messages in return
func SetErrResponse(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	errorMessage, status := "request failed", http.StatusBadRequest

	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		errorMessage, status = "password is invalid", http.StatusForbidden
	case errors.Is(err, cryptography.ErrPasswordTooShort):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, ErrTooBigBody):
		errorMessage, status = err.Error(), http.StatusRequestEntityTooLarge
	case errors.Is(err, cryptography.ErrPasswordTooLong):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, ErrInvalidBody):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, ErrAuthMissing):
		errorMessage, status = err.Error(), http.StatusUnauthorized
	case errors.Is(err, ErrNotFound):
		errorMessage, status = err.Error(), http.StatusNotFound
	case errors.Is(err, sql.ErrNoRows):
		errorMessage, status = err.Error(), http.StatusNoContent
	case errors.Is(err, ErrEmptyResponse):
		errorMessage, status = err.Error(), http.StatusNoContent

	default:
		log.Println("error: ", err.Error())
	}

	ctx.AbortWithStatusJSON(status, gin.H{
		"error": errorMessage,
	})
}
