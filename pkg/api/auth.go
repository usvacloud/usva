package api

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Functions to help with most common tasks
func (s *Server) authorizeRequest(ctx *gin.Context, filename string) bool {
	pwdhash, err := s.db.GetPasswordHash(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return false
	}

	if !pwdhash.Valid {
		return true
	}

	fileauthcookie := fmt.Sprintf("usva-auth-%s", filename)
	authcookieValue, _ := ctx.Cookie(fileauthcookie)

	accesstoken, err := s.db.GetAccessToken(ctx, filename)
	if err != nil {
		setErrResponse(ctx, errAuthFailed)
		return false
	}

	if authcookieValue == accesstoken {
		return true
	}

	pwd, err := s.parseFilePassword(ctx, filename)
	if err != nil {
		setErrResponse(ctx, err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwdhash.String), []byte(pwd))
	if err != nil {
		setErrResponse(ctx, err)
		return false
	}

	ctx.SetCookie(fileauthcookie, accesstoken, s.api.CookieSaveTime, "/", s.api.APIDomain, s.api.UseSecureCookie, true)

	return true
}

func (s *Server) parseFilePassword(ctx *gin.Context, filename string) (string, error) {
	passwordcookie := fmt.Sprintf("usva-password-%s", filename)

	if cookie, err := ctx.Cookie(passwordcookie); err == nil && cookie != "" {
		dec, err := base64.RawStdEncoding.DecodeString(cookie)
		return string(dec), err
	}

	authheader := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
	if len(authheader) < 2 {
		return "", errAuthMissing
	}

	p, err := base64.RawStdEncoding.DecodeString(authheader[1])
	if err != nil {
		return "", err
	}

	ctx.SetCookie(passwordcookie, authheader[1], s.api.CookieSaveTime, "/", s.api.APIDomain, s.api.UseSecureCookie, true)

	return strings.TrimSpace(string(p)), nil
}

func bcryptPasswordHash(pwd []byte) ([]byte, error) {
	pwdlen := len(pwd)
	switch {
	case pwdlen == 0:
		return []byte{}, nil
	case pwdlen > 512:
		return []byte{}, errInvalidBody
	case pwdlen < 6:
		return []byte{}, errInvalidBody
	default:
		return bcrypt.GenerateFromPassword(pwd, 12)
	}
}
