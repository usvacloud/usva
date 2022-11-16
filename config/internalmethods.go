package config

import (
	"log"

	"github.com/romeq/usva/utils"
)

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
