package feedback

import (
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
)

type Handler struct {
	db *db.Queries
}

func NewFeedbackHandler(s *api.Server) *Handler {
	return &Handler{
		db: s.DB,
	}
}
