package arguments

import (
	"flag"

	"github.com/romeq/usva/config"
)

type Arguments struct {
	Config     config.Config
	ConfigFile string
	LogOutput  string
}

func Parse() *Arguments {
	args := &Arguments{}

	// Server configuration
	flag.StringVar(&args.Config.Server.Address, "a", "", "server address")
	flag.IntVar(&args.Config.Server.Port, "p", 0, "server port")

	// File locations
	flag.StringVar(&args.ConfigFile, "c", "/etc/usva/usva.toml", "config location")
	flag.StringVar(&args.LogOutput, "l", "", "log location")

	// Processing arguments
	flag.Parse()
	return args
}
