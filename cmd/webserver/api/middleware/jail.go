package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/pkg/ratelimit"
)

func (s *Handler) Jail(ctx *gin.Context) {
	iphash := ctx.Writer.Header().Get(ratelimit.Headers.Identifier)
	_, err := s.db.IsBanned(ctx.Request.Context(), iphash)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.Next()
		return
	}

	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.AbortWithStatus(http.StatusForbidden)
}
