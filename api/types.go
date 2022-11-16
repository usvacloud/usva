package api

// Limits-struct is primarily used to configure the ratelimiting functionality
type Limits struct {
	Requests int16 // Requests to allow
	Time     int32 // Time to next reset in seconds
}

type Ratelimits struct {
	StrictLimit Limits
	QueryLimit  Limits
}

type APIConfiguration struct {
	MaxSingleUploadSize int64
	MaxUploadSizePerDay int64
	UploadsDir          string
}
