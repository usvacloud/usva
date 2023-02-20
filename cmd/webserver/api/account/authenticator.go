package account

import (
	"context"
	"errors"
	"io"

	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticator struct {
	db *db.Queries
}

func NewAuthenticator(dbconn *db.Queries) UserAuthenticator {
	return UserAuthenticator{
		db: dbconn,
	}
}

func (a UserAuthenticator) NewSession(ctx context.Context, props Login) (string, error) {
	pwd, err := a.db.GetAccountPasswordHash(ctx, props.Username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(props.Password)); err != nil {
		return "", err
	}

	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	return a.db.NewSession(ctx, db.NewSessionParams{
		SessionID: token,
		Username:  props.Username,
	})
}

func (a UserAuthenticator) Authenticate(ctx context.Context, session string) (Account, error) {
	m, err := a.db.GetSessionAccount(ctx, session)
	if errors.Is(err, io.EOF) {
		return Account{}, api.ErrAuthFailed
	}
	return Account(m), nil
}
