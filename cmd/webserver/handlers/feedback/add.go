package feedback

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/generated/db"
)

func (s *Handler) AddFeedback(ctx *gin.Context) {
	body := struct {
		Message string
		Boxes   []int
	}{}
	if err := ctx.BindJSON(&body); err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	maxint := 6
	sort.Ints(body.Boxes)
	if len(body.Boxes) < 1 {
		handlers.SetErrResponse(ctx, handlers.ErrInvalidBody)
		return
	}

	if body.Boxes[len(body.Boxes)-1] > maxint {
		handlers.SetErrResponse(ctx, handlers.ErrInvalidBody)
		return
	}

	var boxestobeinserted string
	for iter := 0; iter < len(body.Boxes); iter++ {
		boxestobeinserted += fmt.Sprint(body.Boxes[iter])
		if iter < len(body.Boxes)-1 {
			boxestobeinserted += ","
		}
	}

	if err := s.db.NewFeedback(ctx, db.NewFeedbackParams{
		Comment: sql.NullString{String: body.Message, Valid: body.Message != ""},
		Boxes:   boxestobeinserted,
	}); err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Feedback added",
	})
}
