package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	errEmptyResponse = errors.New("content was not found")
)

func RestrictionsHandler(ctx *gin.Context, cfg *APIConfiguration) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSize": cfg.MaxSingleUploadSize,
	})
}

func NotFoundHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
	ctx.File("public/404.html")
}
