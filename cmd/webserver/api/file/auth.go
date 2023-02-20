package file

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
	"golang.org/x/crypto/bcrypt"
)

type FileAuthenticator struct {
	db     *db.Queries
	config *api.Configuration
}

func NewFileAuthenticator(db *db.Queries, c *api.Configuration) FileAuthenticator {
	return FileAuthenticator{
		db:     db,
		config: c,
	}
}

// Functions to help with most common tasks
func (a *FileAuthenticator) AuthorizeRequest(ctx *gin.Context, filename string) bool {
	pwdhash, err := a.db.GetPasswordHash(ctx, filename)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return false
	}

	if !pwdhash.Valid {
		return true
	}

	fileauthcookie := fmt.Sprintf("usva-auth-%s", filename)
	authcookieValue, _ := ctx.Cookie(fileauthcookie)

	accesstoken, err := a.db.GetAccessToken(ctx, filename)
	if err != nil {
		api.SetErrResponse(ctx, api.ErrAuthFailed)
		return false
	}

	if authcookieValue == accesstoken {
		return true
	}

	pwd, err := a.ParseFilePassword(ctx, filename)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwdhash.String), []byte(pwd))
	if err != nil {
		api.SetErrResponse(ctx, err)
		return false
	}
	ctx.SetCookie(fileauthcookie, accesstoken, a.config.CookieSaveTime,
		"/", a.config.APIDomain, a.config.UseSecureCookie, true)
	return true
}

func (a *FileAuthenticator) ParseFilePassword(ctx *gin.Context, filename string) (string, error) {
	passwordcookie := fmt.Sprintf("usva-password-%s", filename)

	if cookie, err := ctx.Cookie(passwordcookie); err == nil && cookie != "" {
		dec, err := base64.RawStdEncoding.DecodeString(cookie)
		return string(dec), err
	}

	authheader := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
	if len(authheader) < 2 {
		return "", api.ErrAuthMissing
	}

	p, err := base64.RawStdEncoding.DecodeString(authheader[1])
	if err != nil {
		return "", err
	}

	ctx.SetCookie(passwordcookie, authheader[1], a.config.CookieSaveTime,
		"/", a.config.APIDomain, a.config.UseSecureCookie, true)

	return strings.TrimSpace(string(p)), nil
}

func BCryptPasswordHash(pwd []byte) ([]byte, error) {
	pwdlen := len(pwd)
	switch {
	case pwdlen == 0:
		return []byte{}, nil
	case pwdlen > 512:
		return []byte{}, api.ErrInvalidBody
	case pwdlen < 6:
		return []byte{}, api.ErrInvalidBody
	default:
		return bcrypt.GenerateFromPassword(pwd, 15)
	}
}
