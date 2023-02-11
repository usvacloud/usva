package middleware

import (
	"github.com/romeq/usva/internal/generated/db"
)

type MiddlewareHandler struct {
	db *db.Queries
}

func NewMiddlewareHandler(dbq *db.Queries) *MiddlewareHandler {
	return &MiddlewareHandler{
		db: dbq,
	}
}
