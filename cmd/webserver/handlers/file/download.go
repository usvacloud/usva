package file

import (
	"errors"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/pkg/cryptography"
)

func (s *Handler) DownloadFile(ctx *gin.Context) {
	filename, filenameGiven := ctx.GetQuery("filename")
	if !filenameGiven {
		handlers.SetErrResponse(ctx, errors.New("filename not given"))
		return
	}

	// authorize request
	if !s.auth.AuthorizeRequest(ctx, filename) {
		return
	}

	filepath := path.Join(s.api.UploadsDir, filename)
	fileHandle, err := os.Open(filepath)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	headerPassword, err := s.auth.ParseFilePassword(ctx, filename)
	if err != nil && !errors.Is(err, handlers.ErrAuthMissing) {
		handlers.SetErrResponse(ctx, err)
		return
	}

	encryptionIv, err := s.db.GetEncryptionIV(ctx, filename)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	if len(encryptionIv) == 0 {
		ctx.FileAttachment(filepath, path.Base(filepath))
		return
	} else if errors.Is(err, handlers.ErrAuthMissing) {
		handlers.SetErrResponse(ctx, handlers.ErrAuthMissing)
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", `attachment;`)

	ctx.Status(http.StatusOK)

	derivedKey, err := cryptography.DeriveBasicKey([]byte(headerPassword), s.encryptionKeySize)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	err = cryptography.DecryptStream(ctx.Writer, fileHandle, derivedKey, encryptionIv)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}
}
