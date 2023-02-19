package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
)

func (h Handler) Profile(ctx *gin.Context) {
	s, err := Authenticate(ctx, h.authenticator)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, s)
}
