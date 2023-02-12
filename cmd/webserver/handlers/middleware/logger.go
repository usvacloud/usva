package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Handler) Log(ctx *gin.Context) {
	t := time.Now()
	ctx.Next()
	c := time.Since(t).Milliseconds()

	log.Printf("request: %s %s %d (took %dms) \n",
		ctx.Request.Method, ctx.Request.URL, ctx.Writer.Status(), c)
}
