package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/tapsa/config"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	router.NoRoute(notfound)

	// API
	api := router.Group("/api")
	{
		api.POST("/file/upload", func(ctx *gin.Context) {
			uploadFile(ctx, &cfg.Files)
		})

		api.GET("/file/info", fileInformation)
		api.GET("/file/get", func(ctx *gin.Context) {
			downloadFile(ctx, &cfg.Files)
		})
	}
}

func notfound(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
	ctx.File("public/404.html")
}
