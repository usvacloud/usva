package account

import (
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
)

const (
	accountPasswordBcryptCost = 12
)

type Handler struct {
	dbconn        *db.Queries
	configuration api.Configuration
	authenticator Authenticator
}

func NewAccountsHandler(
	dbconn *db.Queries,
	config api.Configuration,
	authenticator Authenticator,
) *Handler {
	return &Handler{
		dbconn:        dbconn,
		authenticator: authenticator,
	}
}
