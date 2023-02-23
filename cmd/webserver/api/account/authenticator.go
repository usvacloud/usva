package account

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/internal/generated/db"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticator struct {
	db       *db.Queries
	duration time.Duration
}

func NewAuthenticator(dbconn *db.Queries, dur time.Duration) UserAuthenticator {
	return UserAuthenticator{
		db:       dbconn,
		duration: dur,
	}
}

func (a UserAuthenticator) NewSession(ctx context.Context, props Login) (string, error) {
	pwd, err := a.db.GetAccountPasswordHash(ctx, props.Username)
	if err != nil {
		return "", api.ErrAuthFailed
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(props.Password)); err != nil {
		return "", api.ErrAuthFailed
	}

	return a.newsession(ctx, props.Username)
}

func (a UserAuthenticator) Register(ctx context.Context, props db.NewAccountParams) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(props.Password), accountPasswordBcryptCost)
	if err != nil {
		return "", err
	}

	props.Password = string(password)

	ac, err := a.db.NewAccount(ctx, props)
	if err != nil {
		return "", err
	}

	return a.newsession(ctx, ac.Username)
}

func (a UserAuthenticator) Authenticate(ctx context.Context, session string) (Account, error) {
	m, err := a.db.GetSessionAccount(ctx, session)
	if errors.Is(err, pgx.ErrNoRows) {
		return Account{}, api.ErrAuthFailed
	}
	return Account(m), err
}

func (a UserAuthenticator) newsession(ctx context.Context, username string) (string, error) {
	token, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	return a.db.NewAccountSession(ctx, db.NewAccountSessionParams{
		SessionID:  token,
		Username:   username,
		ExpireDate: time.Now().Add(time.Duration(a.duration) * 10),
	})
}

func generateSessionToken() (string, error) {
	randomID := make([]byte, 20)
	if n, err := rand.Read(randomID); err != nil || n < 20 {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomID), nil
}
