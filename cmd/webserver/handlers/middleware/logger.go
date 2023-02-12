package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Handler) Log(ctx *gin.Context) {
	t := time.Now()
	ctx.Next()
	c := time.Since(t).Microseconds()

	log.Printf("'%s %s' %d @ [%.2fms / %dµs]\n",
		ctx.Request.Method, ctx.Request.URL, ctx.Writer.Status(), float64(c)/1000, c)
}
