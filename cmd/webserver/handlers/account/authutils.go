package account

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/generated/db"
)

const sessionTokenCookieName = "session"

// Authenticate parses
func Authenticate(ctx *gin.Context, a Authenticator) (Session, error) {
	token, err := ParseRequestSession(ctx)
	if err != nil {
		return Session{}, err
	}

	m, err := a.AuthenticateSession(ctx.Request.Context(), token)
	if err != nil {
		return Session{}, err
	}

	return Session{
		Token:   token,
		Account: m,
	}, err
}

func (h Handler) newSession(ctx *gin.Context, username string) error {
	sessionToken, err := generateSessionToken()
	if err != nil {
		return err
	}

	session, err := h.dbconn.NewSession(ctx, db.NewSessionParams{
		SessionID: sessionToken,
		Username:  username,
	})
	if err != nil {
		return err
	}

	h.persistSession(ctx, session)
	return nil
}

func (h Handler) persistSession(ctx *gin.Context, token string) {
	ctx.SetCookie(
		sessionTokenCookieName,
		token,
		int(time.Hour)*24,
		"/",
		h.configuration.APIDomain,
		h.configuration.UseSecureCookie,
		true,
	)
}

type Session struct {
	Token   string
	Account db.GetSessionAccountRow
}

func ParseRequestSession(ctx *gin.Context) (string, error) {
	cookie, err := ctx.Cookie(sessionTokenCookieName)
	if cookie == "" || err != nil {
		return "", handlers.ErrAuthMissing
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
