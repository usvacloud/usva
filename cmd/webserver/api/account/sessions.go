package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
)

func (h Handler) Sessions(ctx *gin.Context) {
	token, err := ParseRequestSession(ctx)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	sessionlist, err := h.dbconn.GetSessions(ctx, token)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sessions": sessionlist,
	})
}

func (h Handler) RemoveSessions(ctx *gin.Context) {
	session, err := Authenticate(ctx, h.authenticator)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	m, err := h.dbconn.DeleteSessions(ctx, session.Token)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"removed": m,
	})
}

type removeSessionStruct struct {
	Token string `json:"token"`
}

func (h Handler) RemoveSession(ctx *gin.Context) {
	body, err := api.BindBodyToStruct(ctx, func(rss removeSessionStruct) error {
		if rss.Token == "" {
			return api.ErrInvalidBody
		}
		return nil
	})
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	s, err := Authenticate(ctx, h.authenticator)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return

	}

	if _, err = h.dbconn.DeleteSession(ctx, db.DeleteSessionParams{
		SessionID:   s.Token,
		SessionID_2: body.Token,
	}); err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
