package account

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
)

const sessionTokenCookieName = "session"

type Session struct {
	Token   string
	Account Account
}

func Authenticate(ctx *gin.Context, a auther) (Session, error) {
	token, err := ParseRequestSession(ctx)
	if err != nil {
		return Session{}, err
	}

	m, err := a.Authenticate(ctx.Request.Context(), token)
	if err != nil {
		return Session{}, err
	}

	return Session{
		Token:   token,
		Account: m,
	}, err
}

func (h Handler) persistSession(ctx *gin.Context, token string) {
	ctx.SetCookie(sessionTokenCookieName, token, int(time.Hour)*24, "/",
		h.configuration.APIDomain, h.configuration.UseSecureCookie, true)
}

func (h Handler) newSession(c *gin.Context, u string) (string, error) {
	s, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	return h.dbconn.NewSession(c, db.NewSessionParams{
		SessionID: s,
		Username:  u,
	})
}

func ParseRequestSession(ctx *gin.Context) (string, error) {
	cookie, err := ctx.Cookie(sessionTokenCookieName)
	if cookie == "" || err != nil {
		return "", api.ErrAuthMissing
	}

	return cookie, nil
}

func generateSessionToken() (string, error) {
	randomID := make([]byte, 20)
	if n, err := rand.Read(randomID); err != nil || n < 20 {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomID), nil
}
