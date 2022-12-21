package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/jobscheduler"
	"github.com/romeq/usva/api"
	"github.com/romeq/usva/api/middleware"
	"github.com/romeq/usva/config"
)

func parseRatelimits(cfg *config.Ratelimit) api.Ratelimits {
	return api.Ratelimits{
		StrictLimit: api.Limits(cfg.StrictLimit),
		QueryLimit:  api.Limits(cfg.QueryLimit),
	}
}

func asSeconds(i uint32) time.Duration {
	return time.Duration(i) * time.Second
}

func setupRouteHandlers(router *gin.Engine, cfg *config.Config) {
	// Declare ratelimiters
	strictrl := middleware.NewRatelimiter()
	queryrl := middleware.NewRatelimiter()

	var jobs []jobscheduler.Job
	jobs = append(jobs, jobscheduler.NewJob(0, time.Second, strictrl.Clean, true))
	jobs = append(jobs, jobscheduler.NewJob(0, time.Second, queryrl.Clean, true))

	if cfg.Files.RemoveFilesAfterInactivity {
		job := jobscheduler.NewJob(0, time.Hour*24, func() {
			removeOldFiles(
				asSeconds(cfg.Files.InactivityUntilDelete),
				cfg.Files.UploadsDir,
				cfg.Files.CleanTrashes,
			)
		}, true)

		jobs = append(jobs, job)
	}

	go jobscheduler.Run(jobs)

	api.AuthSaveTime = cfg.Files.AuthSaveTime
	api.AuthUseSecureCookie = cfg.Files.AuthUseSecureCookie
	api.CookieDomain = cfg.Server.ApiDomain
	apic := api.APIConfiguration{
		MaxSingleUploadSize: cfg.Files.MaxSingleUploadSize,
		MaxUploadSizePerDay: cfg.Files.MaxUploadSizePerDay,
		UploadsDir:          cfg.Files.UploadsDir,
	}

	ratelimits := parseRatelimits(&cfg.Ratelimit)

	strict := strictrl.RestrictRequests(ratelimits.StrictLimit.Requests,
		asSeconds(ratelimits.StrictLimit.Time))

	query := queryrl.RestrictRequests(ratelimits.QueryLimit.Requests,
		asSeconds(ratelimits.QueryLimit.Time))

	uploadRestrictor := strictrl.RestrictUploads(time.Duration(24)*time.Hour, cfg.Files.MaxUploadSizePerDay)

	// Middleware/general stuff
	router.Use(middleware.IdentifierHeader)
	router.NoRoute(api.NotFoundHandler)
	router.GET("/restrictions", func(ctx *gin.Context) {
		api.RestrictionsHandler(ctx, &apic)
	})
	router.POST("/",
		strict,
		uploadRestrictor,
		api.UploadFileSimple(strictrl, &apic),
	)

	// Files API
	file := router.Group("/file")
	{
		// Routes
		file.POST("/report", strict, api.ReportFile())
		file.GET("/info", query, api.FileInformation(&apic))
		file.GET("/", query, api.DownloadFile(&apic))
		// TODO: think of a better way to do handling
		//       idea: return random token which is authorized to make changes
		//file.DELETE("/", query, api.DeleteFile(&apic))
		file.POST(
			"/upload",
			strict,
			uploadRestrictor,
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
