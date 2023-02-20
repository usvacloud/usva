package account

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
	"github.com/romeq/usva/internal/generated/db"
)

const (
	maxQueriedFiles = 20
)

func (h Handler) GetOwnedFiles(ctx *gin.Context) {
	session, err := Authenticate(ctx, h.authenticator)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil || limit > maxQueriedFiles {
		limit = maxQueriedFiles
	}

	files, err := h.dbconn.GetSessionOwnerFiles(ctx, db.GetSessionOwnerFilesParams{
		SessionID: session.Token,
		Limit:     int32(limit),
	})
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func (h Handler) GetAllOwnedFiles(ctx *gin.Context) {
	session, err := Authenticate(ctx, h.authenticator)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}
	files, err := h.dbconn.GetAllSessionOwnerFiles(ctx, session.Token)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}
