package api

// API-package is designed to be very modular and independant of other usva's
// packages. It exports the handlers used to handle all existing operations.
// The point of this action was to make using usva as a library possible.

// Limits-struct is primarily used to configure the ratelimiting functionality
type Limits struct {
	AllowedRequests int16 // Requests to allow
	ResetTime       int32 // Time to next reset in seconds
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
