package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/config"
)

func restrictionsHandler(ctx *gin.Context, cfg *config.Config) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSize": cfg.Files.MaxSize,
	})
}

func notFoundHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
	ctx.File("public/404.html")
}
