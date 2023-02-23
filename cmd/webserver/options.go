package main

import (
	"fmt"

	"github.com/usvacloud/usva/cmd/webserver/arguments"
	"github.com/usvacloud/usva/cmd/webserver/config"
	"github.com/usvacloud/usva/internal/utils"
)

type Options config.Config

func NewOptions(cfg *config.Config, args *arguments.Arguments) *Options {
	return &Options{
		Server: config.Server{
			Address: utils.StringOr(args.Config.Server.Address, cfg.Server.Address),
			Port:    utils.IntOr(args.Config.Server.Port, cfg.Server.Port),
		},
	}
}

func (o *Options) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", o.Server.Address, o.Server.Port)
}
