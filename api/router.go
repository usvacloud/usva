package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/api/middleware"
	"github.com/romeq/usva/config"
)

var routeSetup *config.Config

type Limits struct {
	AllowedRequests int
	Time            time.Duration
}
type Ratelimits struct {
	HardLimit Limits
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, ratelimits Ratelimits) {
	routeSetup = cfg
	requestLimiter := middleware.NewRatelimiter()
	hardLimiter := requestLimiter.Limit(
		ratelimits.HardLimit.AllowedRequests,
		ratelimits.HardLimit.Time,
	)

	// Middleware
	router.Use(middleware.IdentifierHeader)
	router.Use(func(ctx *gin.Context) {
		if time.Since(requestLimiter.LastReset) > time.Duration(24)*time.Hour {
			go requestLimiter.Clean()
		}
	})

	// General handlers
	router.NoRoute(notFoundHandler)
	router.GET("/restrictions", func(ctx *gin.Context) {
		restrictionsHandler(ctx, cfg)
	})

	// Files API
	file := router.Group("/file")
	{
		// Routes
		file.GET("/info", initFilesRoute(fileInformation))
		file.GET("/", initFilesRoute(downloadFile))
		file.DELETE("/", initFilesRoute(deleteFile))
		file.POST(
			"/upload",
			requestLimiter.Limit(1, time.Second*time.Duration(2)),
			requestLimiter.LimitDependsBodySize(
				time.Hour*time.Duration(24),
				cfg.Files.MaxUploadSizePerDay,
			),
			uploadFile(requestLimiter, &cfg.Files),
		)
	}

	// Feedback
	feedback := router.Group("/feedback")
	feedback.Use(hardLimiter)
	{
		feedback.GET("/", getFeedback)
		feedback.POST("/", hardLimiter, addFeedback)
	}
}
