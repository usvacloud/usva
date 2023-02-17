package feedback

import (
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/generated/db"
)

type Handler struct {
	db *db.Queries
}

func NewFeedbackHandler(s *handlers.Server) *Handler {
	return &Handler{
		db: s.DB,
	}
}
