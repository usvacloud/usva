package arguments

import (
	"flag"

	"github.com/romeq/usva/internal/config"
)

type Arguments struct {
	Config     config.Config
	ConfigFile string
	LogOutput  string
}

func Parse() *Arguments {
	// Server configuration
	var args Arguments
	flag.StringVar(&args.Config.Server.Address, "a", "", "server address")
	flag.UintVar(&args.Config.Server.Port, "p", 0, "server port")

	// File locations
	flag.StringVar(&args.ConfigFile, "c", "/etc/usva/usva.toml", "config location")
	flag.StringVar(&args.LogOutput, "l", "", "logging location")

	// Processing arguments
	flag.Parse()
	return &args
}
