package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/usvacloud/usva/cmd/webserver/api"
	"net/http/httputil"
)

func (s *Handler) CheckAPIKey(password string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader(api.Headers.APIKeyHeaderName)

		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(ctx.Request, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))

		if apiKey == "" {
			api.SetErrResponse(ctx, api.ErrAPIKeyMissing)
			return
		}
		if apiKey != password {
			api.SetErrResponse(ctx, api.ErrInvalidAPIKey)
			return
		}
	}
}
