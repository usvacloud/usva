package file

import (
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/pkg/authenticator"
)

type Handler struct {
	db                *db.Queries
	config            *api.Configuration
	encryptionKeySize uint32
	auth              authenticator.Authenticator[AuthSession, string, Auth, string]
}

func NewFileHandler(s *api.Server) *Handler {
	return &Handler{
		db:                s.DB,
		config:            s.Config,
		encryptionKeySize: s.EncKeySize,
		auth:              NewAuthenticator(s.DB, s.Config),
	}
}
