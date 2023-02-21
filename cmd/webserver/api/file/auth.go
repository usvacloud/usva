package file

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) authenticate(ctx *gin.Context, filename string) bool {
	password := h.getrequestpassword(ctx, filename)
	sessionCookieName := formatauthcookiename(filename)
	sessionToken, err := ctx.Cookie(sessionCookieName)

	if sessionToken == "" || errors.Is(err, http.ErrNoCookie) {
		session, err := h.auth.NewSession(ctx, NewAuth(filename, password))
		if err != nil {
			api.SetErrResponse(ctx, err)
			return false
		}

		sessionToken = session.Token
	}

	if _, err := h.auth.Authenticate(ctx, NewAuthSession(filename, sessionToken)); err != nil {
		api.SetErrResponse(ctx, err)
		return false
	}

	h.persistSession(ctx, formatauthcookiename(sessionCookieName), sessionToken)

	return true
}

func (h Handler) persistSession(ctx *gin.Context, cookie, value string) {
	ctx.SetCookie(cookie, value, h.config.CookieSaveTime, "/",
		h.config.APIDomain, h.config.UseSecureCookie, true)
}

func (h Handler) getrequestpassword(ctx *gin.Context, filename string) string {
	hdr := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(hdr) == 2 {
		h.persistSession(ctx, formatpasswordcookiename(filename), hdr[1])

		r, _ := base64.RawStdEncoding.DecodeString(hdr[1])
		return strings.TrimSpace(string(r))
	}

	cookie, _ := ctx.Cookie(formatpasswordcookiename(filename))
	r, _ := base64.RawStdEncoding.DecodeString(cookie)
	return strings.TrimSpace(string(r))
}

// passwordhash generates a password hash using constant cost of 12
func (h Handler) passwordhash(s []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(s, 12)
}

// formatauthcookiename returns cookie name where authentication token should be located
func formatauthcookiename(filename string) string {
	return fmt.Sprintf("usva-token-%s", filename)
}

// formatauthcookiename returns cookie name where password should be located
func formatpasswordcookiename(filename string) string {
	return fmt.Sprintf("usva-password-%s", filename)
}
