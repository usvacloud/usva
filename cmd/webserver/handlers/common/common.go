package common

import "github.com/romeq/usva/cmd/webserver/handlers"

type Handler struct {
	config *handlers.Configuration
}

func NewHandler(c *handlers.Configuration) Handler {
	return Handler{
		config: c,
	}
}
