package feedback

import (
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/internal/generated/db"
)

type Handler struct {
	db *db.Queries
}

func NewFeedbackHandler(s *api.Server) *Handler {
	return &Handler{
		db: s.DB,
	}
}
