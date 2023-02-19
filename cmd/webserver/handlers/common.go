package handlers

import "github.com/gin-gonic/gin"

func BindBodyToStruct[T any](ctx *gin.Context, verifyStruct ...func(T) error) (T, error) {
	var body T
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return body, err
	}

	for _, v := range verifyStruct {
		if err := v(body); err != nil {
			return body, err
		}
	}

	return body, nil
}
