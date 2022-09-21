package config

import (
	"io"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/romeq/usva/utils"
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
	AllowedOrigins []string
	TLS            TLS
}
type Files struct {
	MaxSize    int
	UploadsDir string
}
type Database struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}
type Config struct {
	Server   Server
	Files    Files
	Database Database
}

func ParseFromFile(f *os.File) (cfg Config) {
	configtoml, err := io.ReadAll(f)
	utils.Check(err)

	_, err = toml.Decode(string(configtoml), &cfg)
	utils.Check(err)

	cfg.ensureRequiredValues()
	return cfg
}

func (c *Config) ensureRequiredValues() {
	ensureVal("server address", c.Server.Address)
	ensureVal("server port", c.Server.Port)
	ensureVal("file upload directory", c.Files.UploadsDir)
	ensureList("allowed origins", c.Server.AllowedOrigins)
	if c.Server.TLS.Enabled {
		ensureVal("tls: certificate file", c.Server.TLS.CertFile)
		ensureVal("tls: key file", c.Server.TLS.KeyFile)
	}

	c.Database.Host = utils.StringOr(c.Database.Host, "127.0.0.1")
	c.Database.Database = utils.StringOr(c.Database.Database, "usva")
	c.Database.Port = utils.IntOr(c.Database.Port, 5432)
}

func ensureVal(key string, val any) {
	if val == nil || val == "" || val == 0 {
		log.Fatalln(
			"config validation failed: ", key, "is required but not present",
		)
	}
}

func ensureList(key string, val []string) {
	if len(val) == 0 {
		log.Fatalln(
			"config validation failed: list ", key, "is required but not present",
		)
	}
}
