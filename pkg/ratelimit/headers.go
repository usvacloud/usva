package ratelimit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
)

type APIHeaders struct {
	Identifier   string
	AllowedBytes string
}

var Headers = APIHeaders{
	Identifier:   "x-usva-api-identifier",
	AllowedBytes: "x-usva-allowed-bytes",
}

func SetIdentifierHeader(ctx *gin.Context) {
	clientIdentifier := hex.EncodeToString(sha256.New().Sum([]byte(ctx.ClientIP())))
	ctx.Header(Headers.Identifier, clientIdentifier)
	ctx.Next()
}

func setResponseHeaders(ctx *gin.Context, limit, remaining, toreset int16) {
	ctx.Header("RateLimit-Limit", fmt.Sprint(limit))
	ctx.Header("RateLimit-Remaining", fmt.Sprint(remaining))
	ctx.Header("RateLimit-Reset", fmt.Sprint(toreset))
}
