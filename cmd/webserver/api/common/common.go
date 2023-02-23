package common

import "github.com/usvacloud/usva/cmd/webserver/api"

type Handler struct {
	config *api.Configuration
}

func NewHandler(c *api.Configuration) Handler {
	return Handler{
		config: c,
	}
}
