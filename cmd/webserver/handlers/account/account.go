package account

import (
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/generated/db"
)

const (
	accountPasswordBcryptCost = 12
)

type Handler struct {
	dbconn        *db.Queries
	configuration handlers.Configuration
	authenticator Authenticator
}

func NewAccountsHandler(
	dbconn *db.Queries,
	config handlers.Configuration,
	authenticator Authenticator,
) *Handler {
	return &Handler{
		dbconn:        dbconn,
		authenticator: authenticator,
	}
}
