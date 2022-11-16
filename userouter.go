package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/api"
	"github.com/romeq/usva/api/middleware"
	"github.com/romeq/usva/config"
)

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
func parseRatelimits(cfg *config.Ratelimit) api.Ratelimits {
	return api.Ratelimits{
		StrictLimit: api.Limits(cfg.StrictLimit),
		QueryLimit:  api.Limits(cfg.QueryLimit),
	}
}

func SetupRouteHandlers(router *gin.Engine, a *middleware.Ratelimiter, cfg *config.Config) {
	// Declare ratelimiters
	requestLimiter := a
	go initCleanupRoutine(requestLimiter)
	apic := api.APIConfiguration{
		MaxSingleUploadSize: cfg.Files.MaxSingleUploadSize,
		MaxUploadSizePerDay: cfg.Files.MaxUploadSizePerDay,
		UploadsDir:          cfg.Files.UploadsDir,
	}

	ratelimits := parseRatelimits(&cfg.Ratelimit)
	slmt := ratelimits.StrictLimit
	strict := requestLimiter.RestrictRequests(slmt.AllowedRequests, time.Duration(slmt.ResetTime))

	queryLimit := requestLimiter.RestrictRequests(
		ratelimits.QueryLimit.AllowedRequests,
		time.Duration(ratelimits.QueryLimit.ResetTime))

	// Middleware/general stuff
	router.Use(middleware.IdentifierHeader)
	router.NoRoute(api.NotFoundHandler)
	router.GET("/restrictions", func(ctx *gin.Context) {
		api.RestrictionsHandler(ctx, &apic)
	})

	// Files API
	file := router.Group("/file")
	{
		// Routes
		file.GET("/info", queryLimit, api.FileInformation(&apic))
		file.GET("/", queryLimit, api.DownloadFile(&apic))
		file.DELETE("/", queryLimit, api.DeleteFile(&apic))
		file.POST(
			"/upload",
			strict,
			requestLimiter.RestrictUploads(
				time.Duration(24)*time.Hour,
				cfg.Files.MaxUploadSizePerDay,
			),
			api.UploadFile(requestLimiter, &apic),
		)
	}

	// Feedback
	feedback := router.Group("/feedback")
	{
		feedback.GET("/", api.GetFeedback())
		feedback.POST("/", strict, api.AddFeedback())
	}
}
