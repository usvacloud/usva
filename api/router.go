package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/config"
)

var routeSetup *config.Config

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	router.NoRoute(notFoundHandler)
	routeSetup = cfg

	// General handlers
	router.GET("/restrictions", func(ctx *gin.Context) {
		restrictionsHandler(ctx, cfg)
	})

	// Files API
	file := router.Group("/file")
	{
		file.GET("/info", initFilesRoute(fileInformation))
		file.GET("/", initFilesRoute(downloadFile))
		file.DELETE("/", initFilesRoute(deleteFile))
		file.POST("/upload", initFilesRoute(uploadFile))
	}
}
