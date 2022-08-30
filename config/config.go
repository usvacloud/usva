package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/romeq/tapsa/utils"
)

type Server struct {
	Address        string
	Port           int
	TrustedProxies []string
	DebugMode      bool
}
type Files struct {
	MaxSize    int
	UploadsDir string
}
type Config struct {
	Server Server
	Files  Files
}

func New(
	sv_addr string,
	sv_port int,
	sv_tp []string,
	sv_dm bool,
	f_ms int,
	f_ud string,
) Config {
	return Config{
		Server: Server{
			Address:        sv_addr,
			Port:           sv_port,
			TrustedProxies: sv_tp,
			DebugMode:      sv_dm,
		},
		Files: Files{
			MaxSize:    f_ms,
			UploadsDir: f_ud,
		},
	}
}

func ParseFromFile(f *os.File) (cfg Config) {
	configtoml, err := ioutil.ReadAll(f)
	utils.Check(err)

	toml.Decode(string(configtoml), &cfg)
	cfg.EnsureRequiredValues()
	return cfg
}

func (c *Config) EnsureRequiredValues() {
	ensureVal("server address", c.Server.Address)
	ensureVal("server port", c.Server.Port)
	ensureVal("file upload directory", c.Files.UploadsDir)
}

func ensureVal(key string, val any) {
	if val == nil || val == "" || val == 0 {
		log.Fatalln(key, "is required")
	}
}
