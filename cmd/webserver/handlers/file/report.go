package file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/internal/utils"
)

func (s *Handler) ReportFile(ctx *gin.Context) {
	var requestBody struct {
		Filename string `json:"filename"`
		Reason   string `json:"reason"`
	}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	if len(requestBody.Filename) < 36 ||
		!utils.IsBetween(len(requestBody.Reason), 20, 1024) {

		handlers.SetErrResponse(ctx, handlers.ErrInvalidBody)
		return
	}

	err = s.db.NewReport(ctx, db.NewReportParams{
		FileUuid: requestBody.Filename,
		Reason:   requestBody.Reason,
	})
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "thank you! your report has been sent.",
	})
}
