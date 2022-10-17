package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/config"
)

var routeSetup *config.Config

func initFilesRoute(fn func(ctx *gin.Context, files *config.Files)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fn(ctx, &routeSetup.Files)
	}
}

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	routeSetup = cfg

	// Generic handlers
	router.NoRoute(notFoundHandler)
	router.GET("/restrictions", func(ctx *gin.Context) {
		restrictionsHandler(ctx, cfg)
	})

	// Files API
	api := router.Group("/file")
	{
		api.GET("/info", initFilesRoute(fileInformation))
		api.GET("/", initFilesRoute(downloadFile))
		api.DELETE("/", initFilesRoute(deleteFile))
		api.POST("/upload", initFilesRoute(uploadFile))
	}
}
