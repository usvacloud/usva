package middleware

import (
	"github.com/usvacloud/usva/internal/generated/db"
)

type Handler struct {
	db *db.Queries
}

func NewMiddlewareHandler(dbq *db.Queries) *Handler {
	return &Handler{
		db: dbq,
	}
}
