package file

import (
	"context"

	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	db     *db.Queries
	config *api.Configuration
}

func NewAuthenticator(db *db.Queries, c *api.Configuration) Authenticator {
	return Authenticator{
		db:     db,
		config: c,
	}
}

type AuthSession struct {
	filename string
	Token    string
}

func NewAuthSession(f string, t string) AuthSession {
	return AuthSession{
		filename: f,
		Token:    t,
	}
}

func (a Authenticator) Authenticate(c context.Context, auth AuthSession) (string, error) {
	t, err := a.db.GetAccessToken(c, auth.filename)
	if err != nil {
		return "", err
	}

	if t != auth.Token {
		return "", api.ErrAuthFailed
	}

	return t, nil
}

type Auth struct {
	filename string
	password string
}

func NewAuth(file string, password string) Auth {
	return Auth{
		filename: file,
		password: password,
	}
}

func (a Authenticator) NewSession(c context.Context, b Auth) (AuthSession, error) {
	hash, err := a.db.GetPasswordHash(c, b.filename)
	if err != nil {
		return AuthSession{}, err
	}

	if !hash.Valid {
		authSession, err := a.db.GetAccessToken(c, b.filename)
		return NewAuthSession(b.filename, authSession), err
	}

	comparisonError := bcrypt.CompareHashAndPassword([]byte(hash.String), []byte(b.password))
	if comparisonError != nil {
		return AuthSession{}, err
	}

	authSession, err := a.db.GetAccessToken(c, b.filename)
	return NewAuthSession(b.filename, authSession), err
}

func (a Authenticator) Register(c context.Context, b string) (AuthSession, error) {
	authSession, err := a.db.GetAccessToken(c, b)
	return NewAuthSession(b, authSession), err
}
