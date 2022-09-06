package config

import (
	"io"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/romeq/tapsa/utils"
)

type TLS struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}
type Server struct {
	Address        string
	Port           int
	TrustedProxies []string
	DebugMode      bool
	HideRequests   bool
	TLS            TLS
}
type Files struct {
	MaxSize    int
	UploadsDir string
}
type Config struct {
	Server Server
	Files  Files
}

func ParseFromFile(f *os.File) (cfg Config) {
	configtoml, err := io.ReadAll(f)
	utils.Check(err)

	_, err = toml.Decode(string(configtoml), &cfg)
	utils.Check(err)

	cfg.EnsureRequiredValues()
	return cfg
}

func (c *Config) EnsureRequiredValues() {
	ensureVal("server address", c.Server.Address)
	ensureVal("server port", c.Server.Port)
	ensureVal("file upload directory", c.Files.UploadsDir)
	if c.Server.TLS.Enabled {
		ensureVal("tls: certificate file", c.Server.TLS.CertFile)
		ensureVal("tls: key file", c.Server.TLS.KeyFile)
	}
}

func ensureVal(key string, val any) {
	if val == nil || val == "" || val == 0 {
		log.Fatalln("config validation failed: ", key, "is required but not present")
	}
}
