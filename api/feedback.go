package api

import (
	"errors"
	"net/http"
	"regexp"

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

		matched, err := regexp.Match(`1?,2?,3?,4?,5?,6?`, []byte(body.Boxes))
		if !matched || err != nil {
			setErrResponse(ctx, errors.New("invalid body"))
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
		dbFeedbacks, e := dbengine.GetFeedbacks()
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
