package middleware

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func IdentifierHeader(ctx *gin.Context) {
	ctx.Header(
		"Api-Identifier",
		hex.EncodeToString(sha256.New().Sum([]byte(ctx.ClientIP()))),
	)
	ctx.Next()
}
