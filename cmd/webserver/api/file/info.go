package file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/api"
)

func (s *Handler) FileInformation(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Filename not given",
		})
		return
	}

	if !s.authenticate(ctx, filename) {
		return
	}

	f, err := s.db.GetFileInformation(ctx, filename)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	if err = s.db.UpdateLastSeen(ctx, filename); err != nil {
		api.SetErrResponse(ctx, err)
	}

	pwd, err := s.db.GetPasswordHash(ctx, filename)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"filename":   f.FileUuid,
		"size":       f.FileSize,
		"title":      f.Title,
		"uploadDate": f.UploadDate,
		"viewCount":  f.Viewcount,
		"locked":     pwd.Valid,
		"encrypted":  f.Encrypted,
	})
}
