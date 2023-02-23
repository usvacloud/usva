package file

import (
	"errors"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/pkg/cryptography"
)

func (s *Handler) DownloadFile(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		api.SetErrResponse(ctx, errors.New("filename not given"))
		return
	}

	if !s.authenticate(ctx, filename) {
		return
	}

	filepath := path.Join(s.config.UploadsDir, filename)
	fileHandle, err := os.Open(filepath)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	encryptionIv, err := s.db.GetEncryptionIV(ctx, filename)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	if len(encryptionIv) == 0 {
		ctx.FileAttachment(filepath, path.Base(filepath))
		return
	} else if errors.Is(err, api.ErrAuthMissing) {
		api.SetErrResponse(ctx, api.ErrAuthMissing)
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", `attachment;`)

	ctx.Status(http.StatusOK)

	password := s.getrequestpassword(ctx, filename)
	derivedKey, err := cryptography.DeriveBasicKey([]byte(password), s.encryptionKeySize)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	err = cryptography.DecryptStream(ctx.Writer, fileHandle, derivedKey, encryptionIv)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}
}
