package api

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/romeq/tapsa/config"
)

func uploadFile(ctx *gin.Context, uploadOptions *config.Files) {
	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "File upload failed",
		})
		return
	}

	ctx.SaveUploadedFile(f, path.Join(uploadOptions.UploadsDir, f.Filename))
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": "File uploaded",
	})
}

func downloadFile(ctx *gin.Context, downloadOptions *config.Files) {
	ctx.File(path.Join(downloadOptions.UploadsDir, ctx.Param("filename")))
}
