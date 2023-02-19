package common

import "github.com/romeq/usva/cmd/webserver/api"

type Handler struct {
	config *api.Configuration
}

func NewHandler(c *api.Configuration) Handler {
	return Handler{
		config: c,
	}
}
