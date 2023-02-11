package handlers

import "errors"

var (
	ErrAuthMissing = errors.New("missing authentication in a protected file")
	ErrAuthFailed  = errors.New("authorization succeeded but cookies were not set")
	ErrInvalidBody = errors.New("invalid request body")

	ErrNotFound      = errors.New("resource not found")
	ErrEmptyResponse = errors.New("empty response")
	ErrTooBigBody    = errors.New("request body too big")
)
