package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/arguments"
	"github.com/romeq/usva/cmd/webserver/config"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/dbengine"
	"github.com/romeq/usva/internal/utils"
	"github.com/romeq/usva/internal/workers"
)

func main() {
	log.SetFlags(log.Ltime | log.Ldate)

	// arguments
	args := arguments.Parse()
	defer setLogWriter(args.LogOutput).Close()

	// config file
	cfg := config.ParseFromFile(args.ConfigFile)

	// runtime options
	opts := NewOptions(cfg, args)
	db := dbengine.Init(dbengine.DbConfig{
		Port:     cfg.Database.Port,
		Host:     cfg.Database.Host,
		Name:     cfg.Database.Database,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	})

	log.Println("Starting server at", opts.GetListenAddress())

	// start server
	if !cfg.Server.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowCredentials: true,
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Authorization"},
	}))

	utils.Must(r.SetTrustedProxies(cfg.Server.TrustedProxies))

	handler := handlers.NewServer(r, db, &handlers.Configuration{
		MaxEncryptableFileSize: cfg.Files.MaxEncryptableFileSize,
		MaxSingleUploadSize:    cfg.Files.MaxSingleUploadSize,
		MaxUploadSizePerDay:    cfg.Files.MaxUploadSizePerDay,
		UploadsDir:             cfg.Files.UploadsDir,
		CookieSaveTime:         cfg.Files.AuthSaveTime,
		FilePersistDuration:    cfg.Files.InactivityUntilDelete,
		UseSecureCookie:        cfg.Files.AuthUseSecureCookie,
		APIDomain:              cfg.Server.APIDomain,
	}, cfg.Encryption.KeySize)

	trasher := workers.NewTrasher(time.Hour, cfg.Files.InactivityUntilDelete, cfg.Files.UploadsDir, db)
	handler.IncludeServerContextWorker(trasher)

	addRouteHandlers(handler, cfg)

	tlssettings := cfg.Server.TLS
	if tlssettings.Enabled {
		utils.Must(r.RunTLS(opts.GetListenAddress(), tlssettings.CertFile, tlssettings.KeyFile))
	} else {
		utils.Must(r.Run(opts.GetListenAddress()))
	}
}
