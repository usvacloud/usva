package main

import (
	"time"

	"github.com/romeq/usva/cmd/webserver/config"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/cmd/webserver/handlers/account"
	"github.com/romeq/usva/cmd/webserver/handlers/auth"
	"github.com/romeq/usva/cmd/webserver/handlers/common"
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

func addRouteHandlers(server *handlers.Server, cfg *config.Config) {
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
	authhandler := auth.NewAuthHandler(server)

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

	// Files handlers
	fileGroup := router.Group("/file")
	filehandler := file.NewFileHandler(server, authhandler)
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
		accountsGroup.POST("/register", strict, accountsHandler.CreateAccount)
		accountsGroup.POST("/login", strict, accountsHandler.Login)
		accountsGroup.GET("/", query, accountsHandler.Profile)
	}

	sessionsGroup := accountsGroup.Group("/sessions")
	{
		sessionsGroup.GET("/", query, accountsHandler.Sessions)
		sessionsGroup.DELETE("/", query, accountsHandler.RemoveSession)
		sessionsGroup.DELETE("/all", query, accountsHandler.RemoveSessions)
	}
}
