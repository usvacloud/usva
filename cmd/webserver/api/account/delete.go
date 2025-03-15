package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/usvacloud/usva/cmd/webserver/api"
)

func (h Handler) DeleteAccount(ctx *gin.Context) {
	s, err := h.authenticate(ctx, h.authenticator)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return

	}

	if err = h.dbconn.DeleteAccount(ctx, s.Account.AccountID); err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
