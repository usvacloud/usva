package router

import (
	"time"

	"github.com/romeq/usva/cmd/webserver/config"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/cmd/webserver/handlers/auth"
	"github.com/romeq/usva/cmd/webserver/handlers/feedback"
	"github.com/romeq/usva/cmd/webserver/handlers/file"
	"github.com/romeq/usva/cmd/webserver/handlers/middleware"
	"github.com/romeq/usva/internal/workers"
	"github.com/romeq/usva/pkg/ratelimit"
)

func parseRatelimits(cfg *config.Ratelimit) handlers.Ratelimits {
	return handlers.Ratelimits{
		StrictLimit: handlers.Limits(cfg.StrictLimit),
		QueryLimit:  handlers.Limits(cfg.QueryLimit),
	}
}

func ConfigureServer(server *handlers.Server, cfg *config.Config) {
	// Initialize ratelimiters
	strictrl := ratelimit.NewRatelimiter()
	queryrl := ratelimit.NewRatelimiter()

	server.IncludeServerContextWorker(workers.NewRatelimitCleaner(strictrl, time.Second))
	server.IncludeServerContextWorker(workers.NewRatelimitCleaner(queryrl, time.Second))

	ratelimits := parseRatelimits(&cfg.Ratelimit)

	strict := strictrl.RestrictRequests(ratelimits.StrictLimit.Requests, ratelimits.StrictLimit.Time)
	query := queryrl.RestrictRequests(ratelimits.QueryLimit.Requests, ratelimits.QueryLimit.Time)
	uploadRestrictor := strictrl.RestrictUploads(time.Duration(24)*time.Hour, cfg.Files.MaxUploadSizePerDay)

	authhandler := auth.NewAuthHandler(server)
	filehandler := file.NewFileHandler(server, authhandler)
	feedbackhandler := feedback.NewFeedbackHandler(server)
	middlewarehandler := middleware.NewMiddlewareHandler(server.DB)

	// Middleware/general stuff
	router := server.GetRouter()
	router.Use(ratelimit.SetIdentifierHeader)
	router.Use(middlewarehandler.Jail)

	// ungrouped handlers
	{
		router.NoRoute(server.NotFoundHandler)
		router.GET("/restrictions", server.RestrictionsHandler)
	}

	// Files handlers
	fileGroup := router.Group("/file")
	{
		// Routes
		fileGroup.GET("/info", query, filehandler.FileInformation)
		fileGroup.GET("/", query, filehandler.DownloadFile)
		fileGroup.POST("/upload", strict, uploadRestrictor, filehandler.UploadFile)
		fileGroup.POST("/", strict, uploadRestrictor, filehandler.UploadFileSimple)
		fileGroup.POST("/report", strict, filehandler.ReportFile)
	}

	// Feedback
	feedback := router.Group("/feedback")
	{
		feedback.GET("/", query, feedbackhandler.GetFeedback)
		feedback.POST("/", strict, feedbackhandler.AddFeedback)
	}
}
