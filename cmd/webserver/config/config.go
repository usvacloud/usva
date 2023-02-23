package config

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/internal/utils"
)

// ratelimitRestriction includes properties used specifically to configure
// different groups of ratelimits
type Ratelimit struct {
	// StrictLimit struct is used to limit the access of different
	// POST-operations. This works for example in limiting the process
	// of creating a new feedback or as an additional limit to file upload.
	StrictLimit api.Limits `toml:"strict_limit"`

	// QueryLimit configuration is used, by considering it's name, to
	// limit the query operations applied to database etcetera.
	QueryLimit api.Limits `toml:"query_limit"`
}

type Server struct {
	Address        string   `toml:"address"`
	APIDomain      string   `toml:"api_domain"`
	Port           uint     `toml:"port"`
	TrustedProxies []string `toml:"trusted_proxies"`
	DebugMode      bool     `toml:"debug_mode"`
	HideRequests   bool     `toml:"hide_requests"`
	AllowedOrigins []string `toml:"allowed_origins"`
	TLS            struct {
		Enabled  bool   `toml:"enabled"`
		CertFile string `toml:"cert_file"`
		KeyFile  string `toml:"key_file"`
	} `toml:"tls"`
}

type Encryption struct {
	KeySize uint32 `toml:"key_size"`
}

type Files struct {
	MaxEncryptableFileSize     uint64        `toml:"max_encryptable_file_size"`
	MaxSingleUploadSize        uint64        `toml:"max_single_upload_size"`
	MaxUploadSizePerDay        uint64        `toml:"max_upload_size_per_day"`
	UploadsDir                 string        `toml:"uploads_dir"`
	RemoveFilesAfterInactivity bool          `toml:"remove_files_after_inactivity"`
	InactivityUntilDelete      time.Duration `toml:"inactivity_until_delete"`
	AuthSaveTime               time.Duration `toml:"auth_save_time"`
	AuthUseSecureCookie        bool          `toml:"auth_use_secure_cookie"`
}
type Database struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	UseSSL   bool   `toml:"use_ssl"`
}
type Config struct {
	Server     Server     `toml:"server"`
	Files      Files      `toml:"files"`
	Ratelimit  Ratelimit  `toml:"ratelimit"`
	Encryption Encryption `toml:"encryption"`
	Database   Database   `toml:"database"`
}

func ParseFromFile(f string) *Config {
	var cfg Config
	_, err := toml.DecodeFile(f, &cfg)
	utils.Must(err)

	cfg.ensureRequiredValues()
	return &cfg
}
