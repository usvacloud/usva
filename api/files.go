package api

import (
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romeq/tapsa/config"
)

func uploadFile(ctx *gin.Context, uploadOptions *config.Files) {
	f, err := ctx.FormFile("file")
	if err != nil {
		log.Println("error: ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":  "Request failed",
			"reason": err.Error(),
		})
		return
	}
	filename := uuid.New().String() + path.Ext(f.Filename)
	ctx.SaveUploadedFile(f, path.Join(uploadOptions.UploadsDir, filename))

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message":  "File uploaded",
		"filename": filename,
	})
}

func downloadFile(ctx *gin.Context, downloadOptions *config.Files) {
	ctx.File(path.Join(downloadOptions.UploadsDir, ctx.Param("filename")))
}
