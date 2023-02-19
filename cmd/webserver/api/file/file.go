package file

import (
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/cmd/webserver/api/auth"
	"github.com/romeq/usva/internal/generated/db"
)

type Handler struct {
	db                *db.Queries
	api               *api.Configuration
	encryptionKeySize uint32
	auth              *auth.Handler
}

func NewFileHandler(s *api.Server, authHandler *auth.Handler) *Handler {
	return &Handler{
		db:                s.DB,
		api:               s.Config,
		encryptionKeySize: s.EncKeySize,
		auth:              authHandler,
	}
}
