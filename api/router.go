package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/config"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	router.NoRoute(notfound)
	router.GET("/restrictions", func(ctx *gin.Context) {
		restrictions(ctx, cfg)
	})

	// API
	api := router.Group("/file")
	{
		api.GET("/info", fileInformation)
		api.GET("/", func(ctx *gin.Context) {
			downloadFile(ctx, &cfg.Files)
		})

		api.DELETE("/", func(ctx *gin.Context) {
			deleteFile(ctx, &cfg.Files)
		})
		api.POST("/upload", func(ctx *gin.Context) {
			uploadFile(ctx, &cfg.Files)
		})
	}
}

func restrictions(ctx *gin.Context, cfg *config.Config) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSize": cfg.Files.MaxSize,
	})
}

func notfound(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
	ctx.File("public/404.html")
}
