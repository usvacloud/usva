package file

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/cmd/webserver/api/account"
	"github.com/romeq/usva/internal/generated/db"
)

func (h Handler) linkToAccount(ctx *gin.Context, uploadID string) error {
	session, err := account.ParseRequestSession(ctx)
	if errors.Is(err, api.ErrAuthMissing) {
		return nil
	} else if err != nil {
		return err
	}

	return h.db.FileToAccount(ctx, db.FileToAccountParams{
		FileUuid:  uploadID,
		SessionID: session,
	})
}
