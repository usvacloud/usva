package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type createAccountStruct struct {
	Username string
	Password string
}

func verifyCreateProperties(account createAccountStruct) error {
	switch {
	case !utils.IsBetween(len(account.Username), 4, 16):
		return api.ErrUsernameRequirementsNotMet
	case !utils.IsBetween(len(account.Password), 8, 32):
		return api.ErrPasswordRequirementsNotMet
	}

	return nil
}

func (h Handler) CreateAccount(ctx *gin.Context) {
	body, err := api.BindBodyToStruct(ctx, verifyCreateProperties)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), accountPasswordBcryptCost)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	body.Password = string(password)

	if _, err = h.dbconn.NewAccount(ctx.Request.Context(), db.NewAccountParams(body)); err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	s, err := h.newSession(ctx, body.Username)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	h.persistSession(ctx, s)

	ctx.JSON(http.StatusOK, gin.H{
		"username": body.Username,
		"status":   "created",
	})
}