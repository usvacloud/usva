package file

import (
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
)

type Handler struct {
	db                *db.Queries
	config            *api.Configuration
	encryptionKeySize uint32
	auth              FileAuthenticator
}

func NewFileHandler(s *api.Server) *Handler {
	return &Handler{
		db:                s.DB,
		config:            s.Config,
		encryptionKeySize: s.EncKeySize,
		auth:              NewFileAuthenticator(s.DB, s.Config),
	}
}
