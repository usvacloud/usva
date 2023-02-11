package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/pkg/ratelimit"
)

func (s *MiddlewareHandler) Jail(ctx *gin.Context) {
	iphash := ctx.Writer.Header().Get(ratelimit.Headers.Identifier)
	_, err := s.db.IsBanned(ctx.Request.Context(), iphash)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.Next()
		return
	}

	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	} else {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
}
