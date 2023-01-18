package config

import (
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/romeq/usva/internal/api"
	"github.com/romeq/usva/internal/utils"
)

// ratelimitRestriction includes properties used specifically to configure
// different groups of ratelimits
type Ratelimit struct {
	// StrictLimit struct is used to limit the access of different
	// POST-operations. This works for example in limiting the process
	// of creating a new feedback or as an additional limit to file upload.
	StrictLimit api.Limits

	// QueryLimit configuration is used, by considering it's name, to
	// limit the query operations applied to database etcetera.
	QueryLimit api.Limits
}

type Server struct {
	Address        string
	APIDomain      string
	Port           uint
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

type Encryption struct {
	KeySize uint32
}

type Files struct {
	MaxSingleUploadSize        uint64
	MaxUploadSizePerDay        uint64
	UploadsDir                 string
	RemoveFilesAfterInactivity bool
	InactivityUntilDelete      time.Duration
	CleanTrashes               bool
	AuthSaveTime               int
	AuthUseSecureCookie        bool
}

type Config struct {
	Server     Server
	Files      Files
	Ratelimit  Ratelimit
	Encryption Encryption
	Database   struct {
		User     string
		Password string
		Host     string
		Port     int
		Database string
	}
}

func ParseFromFile(f *os.File) *Config {
	configtoml, err := io.ReadAll(f)
	utils.Must(err)

	var cfg Config
	_, err = toml.Decode(string(configtoml), &cfg)
	utils.Must(err)

	cfg.ensureRequiredValues()
	return &cfg
}
