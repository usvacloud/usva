package account

import (
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/internal/generated/db"
	"github.com/usvacloud/usva/pkg/authenticator"
)

const (
	accountPasswordBcryptCost = 12
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type (
	auther  authenticator.Authenticator[string, Account, Login, db.NewAccountParams]
	Account db.GetSessionAccountRow
)

type Handler struct {
	dbconn        *db.Queries
	configuration api.Configuration
	authenticator auther
}

func NewAccountsHandler(
	dbconn *db.Queries,
	config api.Configuration,
	authenticator auther,
) *Handler {
	return &Handler{
		dbconn:        dbconn,
		authenticator: authenticator,
	}
}
