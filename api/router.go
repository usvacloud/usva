package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/api/middleware"
	"github.com/romeq/usva/config"
)

var routeSetup *config.Config

type Limits struct {
	AllowedRequests int16
	Time            time.Duration
}

type Ratelimits struct {
	HardLimit   Limits
	QueryLimit  Limits
	Ratelimiter *middleware.Ratelimiter
}

func SetupRouteHandlers(router *gin.Engine, cfg *config.Config, ratelimits Ratelimits) {
	routeSetup = cfg

	// Declare ratelimiters
	requestLimiter := ratelimits.Ratelimiter
	go initCleanupRoutine(requestLimiter)

	queryLimit := requestLimiter.RestrictRequests(
		ratelimits.QueryLimit.AllowedRequests,
		ratelimits.QueryLimit.Time)

	hardLimiter := requestLimiter.RestrictRequests(
		ratelimits.HardLimit.AllowedRequests,
		ratelimits.HardLimit.Time)

	// Middleware/general stuff
	router.Use(middleware.IdentifierHeader)
	router.NoRoute(notFoundHandler)
	router.GET("/restrictions", func(ctx *gin.Context) {
		restrictionsHandler(ctx, cfg)
	})

	// Files API
	file := router.Group("/file")
	{
		// Routes
		file.GET("/info", queryLimit, initFilesRoute(fileInformation))
		file.GET("/", queryLimit, initFilesRoute(downloadFile))
		file.DELETE("/", queryLimit, initFilesRoute(deleteFile))
		file.POST(
			"/upload",
			hardLimiter,
			requestLimiter.RestrictUploads(
				time.Duration(24)*time.Hour,
				cfg.Files.MaxUploadSizePerDay,
			),
			uploadFile(requestLimiter, &cfg.Files),
		)
	}

	// Feedback
	feedback := router.Group("/feedback")
	{
		feedback.GET("/", getFeedback)
		feedback.POST("/", hardLimiter, addFeedback)
	}
}

func initCleanupRoutine(rt *middleware.Ratelimiter) {
	if rt == nil {
		log.Println("initCleanupRoutine: rt is nil")
		return
	}

	for {
		<-time.After(time.Hour)
		if time.Since(rt.LastCleanup) > time.Duration(24)*time.Hour {
			rt.Clean()
		}
	}
}
