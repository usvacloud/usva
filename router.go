package main

import (
	"time"

	"github.com/romeq/jobscheduler"
	"github.com/romeq/usva/pkg/api"
	"github.com/romeq/usva/pkg/api/middleware"
	"github.com/romeq/usva/pkg/config"
)

func parseRatelimits(cfg *config.Ratelimit) api.Ratelimits {
	return api.Ratelimits{
		StrictLimit: api.Limits(cfg.StrictLimit),
		QueryLimit:  api.Limits(cfg.QueryLimit),
	}
}

func setupRouteHandlers(server *api.Server, cfg *config.Config) {
	// Initialize ratelimiters
	strictrl := middleware.NewRatelimiter()
	queryrl := middleware.NewRatelimiter()

	var jobs []jobscheduler.Job
	jobs = append(jobs, jobscheduler.NewJob(0, time.Second, strictrl.Clean, true))
	jobs = append(jobs, jobscheduler.NewJob(0, time.Second, queryrl.Clean, true))

	if cfg.Files.RemoveFilesAfterInactivity {
		job := jobscheduler.NewJob(0, time.Hour*24, func() {
			server.RemoveOldFilesWorker(cfg.Files.InactivityUntilDelete, cfg.Files.UploadsDir, cfg.Files.CleanTrashes)
		}, true)

		jobs = append(jobs, job)
	}

	go jobscheduler.Run(jobs)

	ratelimits := parseRatelimits(&cfg.Ratelimit)

	strict := strictrl.RestrictRequests(ratelimits.StrictLimit.Requests, ratelimits.StrictLimit.Time)
	query := queryrl.RestrictRequests(ratelimits.QueryLimit.Requests, ratelimits.QueryLimit.Time)
	uploadRestrictor := strictrl.RestrictUploads(time.Duration(24)*time.Hour, cfg.Files.MaxUploadSizePerDay)

	// Middleware/general stuff
	router := server.GetRouter()
	router.Use(middleware.SetIdentifierHeader)
	router.NoRoute(server.NotFoundHandler)
	router.GET("/restrictions", server.RestrictionsHandler)
	router.POST("/", strict, uploadRestrictor, server.UploadFileSimple)

	// Files API
	file := router.Group("/file")
	{
		// Routes
		file.POST("/report", strict, server.ReportFile)
		file.GET("/info", query, server.FileInformation)
		file.GET("/", query, server.DownloadFile)
		file.POST("/upload", strict, uploadRestrictor, server.UploadFile)
	}

	// Feedback
	feedback := router.Group("/feedback")
	{
		feedback.GET("/", query, server.GetFeedback)
		feedback.POST("/", strict, server.AddFeedback)
	}
}
