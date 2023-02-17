package feedback

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
)

func (s *Handler) GetFeedback(ctx *gin.Context) {
	dbFeedbacks, e := s.db.GetFeedbacks(ctx, 10)
	if e != nil {
		handlers.SetErrResponse(ctx, e)
		return
	}

	if len(dbFeedbacks) == 0 {
		handlers.SetErrResponse(ctx, handlers.ErrInvalidBody)
		return
	}

	ctx.JSON(http.StatusOK, dbFeedbacks)
}
