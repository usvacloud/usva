package api

import "time"

// Limits-struct is primarily used to configure the ratelimiting functionality
type Limits struct {
	Requests int16         // Requests to allow
	Time     time.Duration // Time to next reset in seconds
}

type Ratelimits struct {
	StrictLimit Limits
	QueryLimit  Limits
}

type Configuration struct {
	MaxSingleUploadSize uint64
	MaxUploadSizePerDay uint64
	UploadsDir          string
	CookieSaveTime      int
	UseSecureCookie     bool
	APIDomain           string
}
