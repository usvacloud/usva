package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/romeq/tapsa/api"
	"github.com/romeq/tapsa/arguments"
	"github.com/romeq/tapsa/config"
	"github.com/romeq/tapsa/utils"
)

func main() {
	setuplogger()

	// arguments
	args := arguments.Parse()
	defer setLogWriter(args.LogOutput).Close()

	// config file
	cfgHandle, err := os.Open(args.ConfigFile)
	utils.Check(err)
	cfg := config.ParseFromFile(cfgHandle)

	// start server
	log.Println("Starting server at", cfg.GetListenAddress())
	r := setuprouter(cfg)
	utils.Check(r.Run(cfg.GetListenAddress()))
}

func setuprouter(cfg config.Config) *gin.Engine {
	if !cfg.Server.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	api.SetupRoutes(r, &cfg)
	r.SetTrustedProxies(cfg.Server.TrustedProxies)
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
