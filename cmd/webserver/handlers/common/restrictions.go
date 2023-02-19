package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONBytes(bytes uint64) gin.H {
	return gin.H{
		"bytes":     bytes,
		"kilobytes": bytes / 1000,
		"megabytes": bytes / 1000 / 1000,
		"gigabytes": bytes / 1000 / 1000 / 1000,
	}
}

func (h Handler) RestrictionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSingleUploadSize":  JSONBytes(h.config.MaxSingleUploadSize),
		"maxDailyUploadSize":   JSONBytes(h.config.MaxUploadSizePerDay),
		"maxEncryptedFileSize": JSONBytes(h.config.MaxEncryptableFileSize),
		"filePersistDuration": gin.H{
			"seconds": h.config.FilePersistDuration.Seconds(),
			"hours":   h.config.FilePersistDuration.Hours(),
			"days":    h.config.FilePersistDuration.Hours() / 24,
		},
	})
}
