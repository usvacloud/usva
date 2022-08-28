package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/tapsa/config"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	router.NoRoute(notfound)

	// Frontend
	router.GET("/", func(ctx *gin.Context) { ctx.File("public/index.html") })
	router.GET("/style.css", func(ctx *gin.Context) { ctx.File("public/style.css") })

	// API
	api := router.Group("/api")
	{
		api.POST("/file/upload", func(ctx *gin.Context) {
			uploadFile(ctx, &cfg.Files)
		})
		api.GET("/file/get/:filename", func(ctx *gin.Context) {
			downloadFile(ctx, &cfg.Files)
		})
	}
}

func notfound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"error": "Route not found",
	})
}
