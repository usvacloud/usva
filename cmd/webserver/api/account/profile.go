package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
)

func (h Handler) Profile(ctx *gin.Context) {
	s, err := Authenticate(ctx, h.authenticator)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, s)
}
