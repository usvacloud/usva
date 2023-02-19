package api

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/pkg/cryptography"
	"golang.org/x/crypto/bcrypt"
)

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
	case errors.Is(err, sql.ErrNoRows):
		errorMessage, status = err.Error(), http.StatusNoContent
	case errors.Is(err, io.EOF):
		errorMessage, status = "failed to read request content", http.StatusBadRequest

	case errors.Is(err, ErrPasswordRequirementsNotMet):
		errorMessage, status = "password length requirements not met", http.StatusBadRequest

	case errors.Is(err, ErrUsernameRequirementsNotMet):
		errorMessage, status = "password length requirements not met", http.StatusBadRequest

	case errors.Is(err, cryptography.ErrPasswordTooShort):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, cryptography.ErrPasswordTooLong):
		errorMessage, status = err.Error(), http.StatusBadRequest

	case errors.Is(err, ErrTooBigBody):
		errorMessage, status = err.Error(), http.StatusRequestEntityTooLarge
	case errors.Is(err, ErrInvalidBody):
		errorMessage, status = err.Error(), http.StatusBadRequest
	case errors.Is(err, ErrAuthMissing):
		errorMessage, status = err.Error(), http.StatusUnauthorized
	case errors.Is(err, ErrNotFound):
		errorMessage, status = err.Error(), http.StatusNotFound
	case errors.Is(err, ErrEmptyResponse):
		errorMessage, status = err.Error(), http.StatusNoContent

	default:
		log.Println("error: ", err.Error())
	}

	ctx.AbortWithStatusJSON(status, gin.H{
		"error": errorMessage,
	})
}
