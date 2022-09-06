package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/tapsa/api"
	"github.com/romeq/tapsa/arguments"
	"github.com/romeq/tapsa/config"
	"github.com/romeq/tapsa/dbengine"
	"github.com/romeq/tapsa/utils"
)

type Options struct {
	server config.Server
}

func parseOpts(cfg config.Config, args arguments.Arguments) Options {
	return Options{
		server: config.Server{
			Address: utils.StringOr(args.Address, cfg.Server.Address),
			Port:    utils.IntOr(args.Port, cfg.Server.Port),
		},
	}
}

func (o *Options) getaddr() string {
	return fmt.Sprintf("%s:%d", o.server.Address, o.server.Port)
}

func setuprouter(cfg config.Config) *gin.Engine {
	if !cfg.Server.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	if !cfg.Server.HideRequests {
		r.Use(requestLogger)
	}

	api.SetupRoutes(r, &cfg)
	utils.Check(r.SetTrustedProxies(cfg.Server.TrustedProxies))

	return r
}

func setuplogger() {
	log.SetFlags(log.Ltime | log.Ldate | log.Lshortfile)
}

func setLogWriter(file string) *os.File {
	if file == "" {
		return nil
	}

	fhandle, err := os.OpenFile(file, os.O_WRONLY, 0644)
	utils.Check(err)

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
	setuplogger()

	// arguments
	args := arguments.Parse()
	defer setLogWriter(args.LogOutput).Close()

	// config file
	cfgHandle, err := os.Open(args.ConfigFile)
	utils.Check(err)
	cfg := config.ParseFromFile(cfgHandle)

	// runtime options
	opts := parseOpts(cfg, args)
	dbengine.Init(args.DatabasePath)
	defer dbengine.DbConnection.Close()

	// start server
	log.Println("Starting server at", opts.getaddr())

	r := setuprouter(cfg)
	tlssettings := cfg.Server.TLS
	if tlssettings.Enabled {
		utils.Check(r.RunTLS(opts.getaddr(), tlssettings.KeyFile, tlssettings.CertFile))
	} else {
		utils.Check(r.Run(opts.getaddr()))
	}
}
