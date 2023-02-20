package main

import (
	"time"

	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/cmd/webserver/api/account"
	"github.com/romeq/usva/cmd/webserver/api/common"
	"github.com/romeq/usva/cmd/webserver/api/feedback"
	"github.com/romeq/usva/cmd/webserver/api/file"
	"github.com/romeq/usva/cmd/webserver/api/middleware"
	"github.com/romeq/usva/cmd/webserver/config"
	"github.com/romeq/usva/internal/workers"
	"github.com/romeq/usva/pkg/ratelimit"
)

func parseRatelimits(cfg *config.Ratelimit) api.Ratelimits {
	return api.Ratelimits{
		StrictLimit: api.Limits(cfg.StrictLimit),
		QueryLimit:  api.Limits(cfg.QueryLimit),
	}
}

func addRouteapi(server *api.Server, cfg *config.Config) {
	// Initialize ratelimiters
	strictrl := ratelimit.NewRatelimiter()
	queryrl := ratelimit.NewRatelimiter()

	server.IncludeServerContextWorker(workers.NewRatelimitCleaner(strictrl, time.Second))
	server.IncludeServerContextWorker(workers.NewRatelimitCleaner(queryrl, time.Second))

	ratelimits := parseRatelimits(&cfg.Ratelimit)

	strict := strictrl.RestrictRequests(ratelimits.StrictLimit.Requests, ratelimits.StrictLimit.Time)
	query := queryrl.RestrictRequests(ratelimits.QueryLimit.Requests, ratelimits.QueryLimit.Time)
	uploadRestrictor := strictrl.RestrictUploads(time.Duration(24)*time.Hour, cfg.Files.MaxUploadSizePerDay)

	// Middleware/general stuff
	router := server.GetRouter()

	// Middlewares
	middlewarehandler := middleware.NewMiddlewareHandler(server.DB)
	{
		router.Use(ratelimit.SetIdentifierHeader)
		router.Use(middlewarehandler.Jail)
		router.NoRoute(server.NotFoundHandler)

		if !cfg.Server.HideRequests {
			router.Use(middlewarehandler.Log)
		}
	}

	// Common
	commonHandler := common.NewHandler(server.Config)
	{
		router.GET("/restrictions", commonHandler.RestrictionsHandler)
	}

	// Files api
	fileGroup := router.Group("/file")
	filehandler := file.NewFileHandler(server)
	{
		// Routes
		fileGroup.GET("/info", query, filehandler.FileInformation)
		fileGroup.GET("/", query, filehandler.DownloadFile)
		fileGroup.POST("/upload", strict, uploadRestrictor, filehandler.UploadFile)
		fileGroup.POST("/report", strict, filehandler.ReportFile)
	}

	// Feedback
	feedbackGroup := router.Group("/feedback")
	feedbackhandler := feedback.NewFeedbackHandler(server)
	{
		feedbackGroup.GET("/", query, feedbackhandler.GetFeedback)
		feedbackGroup.POST("/", strict, feedbackhandler.AddFeedback)
	}

	// Accounts
	accountsGroup := router.Group("/account")
	userAuthenticator := account.NewAuthenticator(server.DB)
	accountsHandler := account.NewAccountsHandler(server.DB, *server.Config, userAuthenticator)
	{
		accountsGroup.GET("/", query, accountsHandler.Profile)
		accountsGroup.GET("/files", query, accountsHandler.GetOwnedFiles)
		accountsGroup.GET("/files/all", query, accountsHandler.GetAllOwnedFiles)
		accountsGroup.POST("/login", strict, accountsHandler.Login)
		accountsGroup.POST("/register", strict, accountsHandler.CreateAccount)
	}

	sessionsGroup := accountsGroup.Group("/sessions")
	{
		sessionsGroup.GET("/", query, accountsHandler.Sessions)
		sessionsGroup.DELETE("/", query, accountsHandler.RemoveSession)
		sessionsGroup.DELETE("/all", query, accountsHandler.RemoveSessions)
	}
}
