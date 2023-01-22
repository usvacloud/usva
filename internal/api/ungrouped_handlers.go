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
	errEmptyResponse = errNotFound
)

func (s *Server) RestrictionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSize": s.api.MaxSingleUploadSize,
	})
}

func (s *Server) NotFoundHandler(ctx *gin.Context) {
	setErrResponse(ctx, errEmptyResponse)
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
	case errors.Is(err, cryptography.ErrPasswordTooLong):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, errInvalidBody):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, errAuthMissing):
		errorMessage, status = err.Error(), http.StatusUnauthorized
	case errors.Is(err, errNotFound):
		errorMessage, status = err.Error(), http.StatusNotFound
	case errors.Is(err, sql.ErrNoRows):
		errorMessage, status = errNotFound.Error(), http.StatusNotFound
	case errors.Is(err, errEmptyResponse):
		errorMessage, status = err.Error(), http.StatusNoContent

	default:
		log.Println("error: ", err.Error())
	}

	ctx.AbortWithStatusJSON(status, gin.H{
		"error": errorMessage,
	})
}