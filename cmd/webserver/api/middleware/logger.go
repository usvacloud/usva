package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/usvacloud/usva/pkg/ratelimit"
)

func (s *Handler) Log(ctx *gin.Context) {
	t := time.Now()
	ctx.Next()
	c := time.Since(t).Microseconds()

	ip := ctx.Writer.Header().Get(ratelimit.Headers.Identifier)
	log.Printf("'%s %s' %d from %s [%.2fms / %dµs]\n",
		ctx.Request.Method, ctx.Request.URL, ctx.Writer.Status(), ip, float64(c)/1000, c)
}
