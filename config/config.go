package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/romeq/usva/api"
	"github.com/romeq/usva/utils"
)

type Server struct {
	Address        string
	Port           int
	TrustedProxies []string
	DebugMode      bool
	HideRequests   bool
	AllowedOrigins []string
	TLS            struct {
		Enabled  bool
		CertFile string
		KeyFile  string
	}
}

// ratelimitRestriction includes properties used specifically to configure
// different groups of ratelimits
type ratelimitRestriction api.Limits

type Ratelimit struct {
	// StrictLimit struct is used to limit the access of different
	// POST-operations. This works for example in limiting the process
	// of creating a new feedback or as an additional limit to file upload.
	StrictLimit ratelimitRestriction

	// QueryLimit configuration is used, by considering it's name, to
	// limit the query operations applied to database etcetera.
	QueryLimit ratelimitRestriction
}

type Files struct {
	MaxSingleUploadSize int64
	MaxUploadSizePerDay int64
	UploadsDir          string
}

type Config struct {
	Server    Server
	Files     Files
	Ratelimit Ratelimit
	Database  struct {
		User     string
		Password string
		Host     string
		Port     int
		Database string
	}
}

func ParseFromFile(f *os.File) (cfg Config) {
	configtoml, err := io.ReadAll(f)
	utils.Check(err)

	_, err = toml.Decode(string(configtoml), &cfg)
	utils.Check(err)

	cfg.ensureRequiredValues()
	return cfg
}
