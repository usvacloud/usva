package api

import (
	"fmt"
)

var (
	ErrAuthMissing                = fmt.Errorf("missing authentication to endpoint even though it's required")
	ErrAuthFailed                 = fmt.Errorf("invalid authentication parameters")
	ErrInvalidBody                = fmt.Errorf("invalid request body")
	ErrNotFound                   = fmt.Errorf("resource was not found")
	ErrEmptyResponse              = fmt.Errorf("empty response")
	ErrTooBigBody                 = fmt.Errorf("request body too big")
	ErrPasswordRequirementsNotMet = fmt.Errorf("password has to be between 8 and 32 charecters long")
	ErrUsernameRequirementsNotMet = fmt.Errorf("username has to be between 4 and 16 charecters long")
	ErrAPIKeyMissing              = fmt.Errorf("api key is missing")
	ErrInvalidAPIKey              = fmt.Errorf("api key is invalid")
)
