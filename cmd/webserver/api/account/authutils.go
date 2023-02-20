package account

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
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

func ParseRequestSession(ctx *gin.Context) (string, error) {
	cookie, err := ctx.Cookie(sessionTokenCookieName)
	if cookie == "" || err != nil {
		return "", api.ErrAuthMissing
	}

	return cookie, nil
}
