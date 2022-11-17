package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/api"
	"github.com/romeq/usva/api/middleware"
	"github.com/romeq/usva/config"
)

func initCleanupRoutine(rt ...*middleware.Ratelimiter) {
	if rt == nil {
		log.Println("initCleanupRoutine: rt is nil")
		return
	}

	for {
		<-time.After(time.Hour)
		for _, ft := range rt {
			if time.Since(ft.LastCleanup) > time.Duration(24)*time.Hour {
				ft.Clean()
			}
		}
	}
}
func parseRatelimits(cfg *config.Ratelimit) api.Ratelimits {
	return api.Ratelimits{
		StrictLimit: api.Limits(cfg.StrictLimit),
		QueryLimit:  api.Limits(cfg.QueryLimit),
	}
}

func toseconds(i uint32) time.Duration {
	return time.Duration(i) * time.Second
}

func setupRouteHandlers(router *gin.Engine, cfg *config.Config) {
	// Declare ratelimiters
	strictrl := middleware.NewRatelimiter()
	queryrl := middleware.NewRatelimiter()
	go initCleanupRoutine(strictrl, queryrl)
	go removeOldFilesWorker(toseconds(cfg.Files.InactivityUntilDelete), cfg.Files.UploadsDir, cfg.Files.CleanTrashes)

	apic := api.APIConfiguration{
		MaxSingleUploadSize: cfg.Files.MaxSingleUploadSize,
		MaxUploadSizePerDay: cfg.Files.MaxUploadSizePerDay,
		UploadsDir:          cfg.Files.UploadsDir,
	}

	ratelimits := parseRatelimits(&cfg.Ratelimit)

	strict := strictrl.RestrictRequests(ratelimits.StrictLimit.Requests,
		toseconds(ratelimits.StrictLimit.Time))

	query := queryrl.RestrictRequests(ratelimits.QueryLimit.Requests,
		toseconds(ratelimits.QueryLimit.Time))

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
		file.POST("/report", strict, api.ReportFile())
		file.GET("/info", query, api.FileInformation(&apic))
		file.GET("/", query, api.DownloadFile(&apic))
		file.DELETE("/", query, api.DeleteFile(&apic))
		file.POST(
			"/upload",
			strict,
			strictrl.RestrictUploads(
				time.Duration(24)*time.Hour,
				cfg.Files.MaxUploadSizePerDay,
			),
			// pass strictrl for updating the uploaded content
			api.UploadFile(strictrl, &apic),
		)
	}

	// Feedback
	feedback := router.Group("/feedback")
	{
		feedback.GET("/", query, api.GetFeedback())
		feedback.POST("/", strict, api.AddFeedback())
	}
}
