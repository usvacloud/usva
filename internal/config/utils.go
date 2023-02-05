package config

import (
	"log"

	"github.com/romeq/usva/internal/utils"
)

func (c *Config) ensureRequiredValues() {
	ensureVal("server: address", c.Server.Address)
	ensureVal("server: port", c.Server.Port)
	ensureList("server: allowed origins", c.Server.AllowedOrigins)
	if c.Server.TLS.Enabled {
		ensureVal("server: tls: certificate file", c.Server.TLS.CertFile)
		ensureVal("server: tls: key file", c.Server.TLS.KeyFile)
	}

	ensureVal("files: upload directory", c.Files.UploadsDir)

	if c.Encryption.KeySize != 16 && c.Encryption.KeySize != 32 && c.Encryption.KeySize != 24 {
		log.Fatalf("encryption: key size is invalid (%d is not 16, 24 or 32).", c.Encryption.KeySize)
	}

	c.Database.Host = utils.StringOr(c.Database.Host, "127.0.0.1")
	c.Database.Database = utils.StringOr(c.Database.Database, "usva")
	c.Database.Port = int(utils.IntOr(uint(c.Database.Port), 5432))
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
