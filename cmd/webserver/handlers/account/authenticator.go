package account

import (
	"context"
	"errors"
	"io"

	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/generated/db"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	db *db.Queries
}

func NewAuthenticator(dbconn *db.Queries) Authenticator {
	return Authenticator{
		db: dbconn,
	}
}

func (a Authenticator) NewSession(ctx context.Context, username string, password string) (string, error) {
	pwd, err := a.db.GetAccountPasswordHash(ctx, username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password)); err != nil {
		return "", err
	}

	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	return a.db.NewSession(ctx, db.NewSessionParams{
		SessionID: token,
		Username:  username,
	})
}

func (a Authenticator) AuthenticateSession(ctx context.Context, session string) (db.GetSessionAccountRow, error) {
	m, err := a.db.GetSessionAccount(ctx, session)
	if errors.Is(err, io.EOF) {
		return db.GetSessionAccountRow{}, handlers.ErrAuthFailed
	}
	return m, nil
}
