package file

import (
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/cmd/webserver/handlers/auth"
	"github.com/romeq/usva/internal/generated/db"
)

type Handler struct {
	db                *db.Queries
	api               *handlers.Configuration
	encryptionKeySize uint32
	auth              *auth.Handler
}

func NewFileHandler(s *handlers.Server, authHandler *auth.Handler) *Handler {
	return &Handler{
		db:                s.DB,
		api:               s.Config,
		encryptionKeySize: s.EncKeySize,
		auth:              authHandler,
	}
}
