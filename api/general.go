package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/config"
)

var (
	errEmptyResponse = errors.New("content was not found")
)

func restrictionsHandler(ctx *gin.Context, cfg *config.Config) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSize": cfg.Files.MaxSingleUploadSize,
	})
}

func notFoundHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
	ctx.File("public/404.html")
}
