package file

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/internal/utils"
)

func (s *Handler) FileInformation(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Filename not given",
		})
		return
	}

	if !s.auth.AuthorizeRequest(ctx, filename) {
		return
	}

	f, err := s.db.FileInformation(ctx, filename)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	if err = s.db.UpdateLastSeen(ctx, filename); err != nil {
		handlers.SetErrResponse(ctx, err)
	}

	pwd, err := s.db.GetPasswordHash(ctx, filename)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	filesize, err := utils.FileSize(path.Join(s.api.UploadsDir, filename))
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"filename":   f.FileUuid,
		"size":       filesize,
		"title":      f.Title,
		"uploadDate": f.UploadDate,
		"viewCount":  f.Viewcount,
		"locked":     pwd.Valid,
		"encrypted":  f.Encrypted,
	})
}
