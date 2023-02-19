package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/cmd/webserver/arguments"
	"github.com/romeq/usva/cmd/webserver/config"
	"github.com/romeq/usva/internal/dbengine"
	"github.com/romeq/usva/internal/utils"
	"github.com/romeq/usva/internal/workers"
)

func main() {
	log.SetFlags(log.Ltime | log.Ldate)

	// arguments
	args := arguments.Parse()
	logWriterHandle := setLogWriter(args.LogOutput)

	// config file
	cfg := config.ParseFromFile(args.ConfigFile)

	// runtime options
	opts := NewOptions(cfg, args)
	db, dbClose := dbengine.Init(dbengine.DbConfig{
		Port:     cfg.Database.Port,
		Host:     cfg.Database.Host,
		Name:     cfg.Database.Database,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		UseSSL:   cfg.Database.UseSSL,
	})

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

	handler := api.NewServer(r, db, &api.Configuration{
		MaxEncryptableFileSize: cfg.Files.MaxEncryptableFileSize,
		MaxSingleUploadSize:    cfg.Files.MaxSingleUploadSize,
		MaxUploadSizePerDay:    cfg.Files.MaxUploadSizePerDay,
		UploadsDir:             cfg.Files.UploadsDir,
		CookieSaveTime:         int(cfg.Files.AuthSaveTime.Seconds()),
		FilePersistDuration:    cfg.Files.InactivityUntilDelete,
		UseSecureCookie:        cfg.Files.AuthUseSecureCookie,
		APIDomain:              cfg.Server.APIDomain,
	}, cfg.Encryption.KeySize)

	trasher := workers.NewTrasher(time.Hour, cfg.Files.InactivityUntilDelete, cfg.Files.UploadsDir, db)
	handler.IncludeServerContextWorker(trasher)

	addRouteapi(handler, cfg)
	srv := http.Server{
		ReadTimeout:       time.Minute,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		Handler:           r,
		Addr:              opts.GetListenAddress(),
	}

	// start server
	log.Printf("starting server at %s (pid %d)", opts.GetListenAddress(), os.Getpid())
	go func() {
		tlssettings := cfg.Server.TLS
		if tlssettings.Enabled {
			err := srv.ListenAndServeTLS(tlssettings.CertFile, tlssettings.KeyFile)
			utils.Must(err)
		} else {
			err := srv.ListenAndServe()
			utils.Must(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	dbClose()
	logWriterHandle.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.Fatal("server forced to shutdown:", err)
	}
}
