package arguments

import (
	"flag"
)

type Arguments struct {
	Address      string
	Port         int
	ConfigFile   string
	LogOutput    string
	DatabasePath string
}

func Parse() (args Arguments) {
	flag.StringVar(&args.Address, "a", "127.0.0.1", "server address")
	flag.StringVar(&args.ConfigFile, "c", "/etc/tapsa/tapsa.toml", "config location")
	flag.StringVar(&args.LogOutput, "l", "", "log location")
	flag.StringVar(&args.DatabasePath, "d", "/usr/share/tapsa/files.db", "database location")
	flag.IntVar(&args.Port, "p", 8080, "server port")

	flag.Parse()
	return args
}
