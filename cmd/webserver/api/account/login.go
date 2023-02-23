package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/internal/utils"
)

func verifyLoginProperties(account Login) error {
	switch {
	case !utils.IsBetween(len(account.Username), 4, 16):
		return api.ErrUsernameRequirementsNotMet
	case !utils.IsBetween(len(account.Password), 8, 32):
		return api.ErrPasswordRequirementsNotMet
	}

	return nil
}

func (h Handler) Login(ctx *gin.Context) {
	b, err := api.BindBodyToStruct(ctx, verifyLoginProperties)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	sessionID, err := h.authenticator.NewSession(ctx, b)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	h.persistSession(ctx, sessionID)

	ctx.JSON(http.StatusOK, gin.H{
		"sessionId": sessionID,
	})
}
