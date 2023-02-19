package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/utils"
)

type loginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func verifyLoginProperties(account loginStruct) error {
	switch {
	case !utils.IsBetween(len(account.Username), 4, 16):
		return handlers.ErrUsernameRequirementsNotMet
	case !utils.IsBetween(len(account.Password), 8, 32):
		return handlers.ErrPasswordRequirementsNotMet
	}

	return nil
}

func (h Handler) Login(ctx *gin.Context) {
	b, err := handlers.BindBodyToStruct(ctx, verifyLoginProperties)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	sessionID, err := h.authenticator.NewSession(ctx, b.Username, b.Password)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	h.persistSession(ctx, sessionID)

	ctx.JSON(http.StatusOK, gin.H{
		"sessionId": sessionID,
	})
}
