package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/config"
)

var routeSetup *config.Config

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	routeSetup = cfg

	// General handlers
	router.NoRoute(notFoundHandler)
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

	// Feedback
	feedback := router.Group("/feedback")
	{
		feedback.GET("/", getFeedback)
		feedback.POST("/", addFeedback)
	}
}
