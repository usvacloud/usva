package feedback

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/usvacloud/usva/cmd/webserver/api"
)

func (s *Handler) GetFeedback(ctx *gin.Context) {
	dbFeedbacks, e := s.db.GetFeedbacks(ctx, 10)
	if e != nil {
		api.SetErrResponse(ctx, e)
		return
	}

	if len(dbFeedbacks) == 0 {
		api.SetErrResponse(ctx, api.ErrInvalidBody)
		return
	}

	ctx.JSON(http.StatusOK, dbFeedbacks)
}
