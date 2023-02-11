package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/arguments"
	"github.com/romeq/usva/cmd/webserver/config"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/cmd/webserver/router"
	"github.com/romeq/usva/internal/dbengine"
	"github.com/romeq/usva/internal/utils"
	"github.com/romeq/usva/internal/workers"
)

// options is a struct that helps with overriding
// configuration values with command line arguments
type Options config.Config

func NewOptions(cfg *config.Config, args *arguments.Arguments) *Options {
	return &Options{
		Server: config.Server{
			Address: utils.StringOr(args.Config.Server.Address, cfg.Server.Address),
			Port:    utils.IntOr(args.Config.Server.Port, cfg.Server.Port),
		},
	}
}

func (o *Options) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", o.Server.Address, o.Server.Port)
}

func setupEngine(cfg *config.Config) *gin.Engine {
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
	if !cfg.Server.HideRequests {
		r.Use(requestLogger)
	}

	utils.Must(r.SetTrustedProxies(cfg.Server.TrustedProxies))
	return r
}

func setLogWriter(file string) *os.File {
	if file == "" {
		return nil
	}

	fhandle, err := os.OpenFile(file, os.O_WRONLY, 0o644)
	utils.Must(err)

	log.SetOutput(fhandle)
	return fhandle
}

func requestLogger(ctx *gin.Context) {
	t := time.Now()
	ctx.Next()
	c := time.Since(t).Milliseconds()

	log.Printf("request: %s %s %d (took %dms) \n",
		ctx.Request.Method, ctx.Request.URL, ctx.Writer.Status(), c)
}

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
	r := setupEngine(cfg)
	server := handlers.NewServer(r, db, &handlers.Configuration{
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
	server.IncludeServerContextWorker(trasher)

	router.ConfigureServer(server, cfg)

	tlssettings := cfg.Server.TLS
	if tlssettings.Enabled {
		utils.Must(r.RunTLS(opts.GetListenAddress(), tlssettings.CertFile, tlssettings.KeyFile))
	} else {
		utils.Must(r.Run(opts.GetListenAddress()))
	}
}
