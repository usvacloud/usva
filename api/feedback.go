package api

import (
    "sort"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/dbengine"
)

func AddFeedback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := dbengine.Feedback{}
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

		if err := dbengine.AddFeedback(&body); err != nil {
			setErrResponse(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Feedback added",
		})
	}
}

func GetFeedback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbFeedbacks, e := dbengine.GetFeedbacks(10)
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
}
