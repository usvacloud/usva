package file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/internal/generated/db"
	"github.com/usvacloud/usva/internal/utils"
)

func (s *Handler) ReportFile(ctx *gin.Context) {
	var requestBody struct {
		Filename string `json:"filename"`
		Reason   string `json:"reason"`
	}
	err := ctx.BindJSON(&requestBody)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	if len(requestBody.Filename) < 36 ||
		!utils.IsBetween(len(requestBody.Reason), 20, 1024) {

		api.SetErrResponse(ctx, api.ErrInvalidBody)
		return
	}

	err = s.db.NewReport(ctx, db.NewReportParams{
		FileUuid: requestBody.Filename,
		Reason:   requestBody.Reason,
	})
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "thank you! your report has been sent.",
	})
}
