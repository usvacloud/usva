package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/db"
	"github.com/romeq/usva/dbengine"
)

func (s *Server) AddFeedback(ctx *gin.Context) {
	body := struct {
		Message string
		Boxes   []int
	}{}
	if err := ctx.BindJSON(&body); err != nil {
		setErrResponse(ctx, err)
		return
	}

	maxint := 6
	sort.Ints(body.Boxes)
	if len(body.Boxes) < 1 {
		setErrResponse(ctx, errInvalidBody)
		return
	}

	if body.Boxes[len(body.Boxes)-1] > maxint {
		setErrResponse(ctx, errInvalidBody)
		return
	}

	var boxestobeinserted string
	for iter := 0; iter < len(body.Boxes); iter++ {
		boxestobeinserted += fmt.Sprint(body.Boxes[iter])
		if iter < len(body.Boxes)-1 {
			boxestobeinserted += ","
		}
	}

	if err := dbengine.DB.NewFeedback(ctx, db.NewFeedbackParams{
		Comment: sql.NullString{String: body.Message, Valid: body.Message != ""},
		Boxes:   boxestobeinserted,
	}); err != nil {
		setErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Feedback added",
	})
}

func (s *Server) GetFeedback(ctx *gin.Context) {
	dbFeedbacks, e := dbengine.DB.GetFeedbacks(ctx, 10)
	if e != nil {
		setErrResponse(ctx, e)
		return
	}

	if len(dbFeedbacks) <= 0 {
		setErrResponse(ctx, errEmptyResponse)
		return
	}

	ctx.JSON(http.StatusOK, dbFeedbacks)
}
