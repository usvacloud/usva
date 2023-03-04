package file

import (
	"database/sql"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/usvacloud/usva/cmd/webserver/api"
	"github.com/usvacloud/usva/internal/generated/db"
	"github.com/usvacloud/usva/pkg/cryptography"
)

// UploadFile
func (s *Handler) UploadFile(ctx *gin.Context) {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	filename := uuid.NewString() + path.Ext(formFile.Filename)

	if s.config.MaxSingleUploadSize > 0 && uint64(formFile.Size) > s.config.MaxSingleUploadSize {
		api.SetErrResponse(ctx, api.ErrTooBigBody)
		return
	}

	formpwd := strings.TrimSpace(ctx.PostForm("password"))
	pwd, err := base64.RawStdEncoding.DecodeString(formpwd)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}
	password := strings.TrimSpace(string(pwd))

	formFileHandle, err := formFile.Open()
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	absUploads, err := filepath.Abs(s.config.UploadsDir)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	file, err := os.OpenFile(path.Join(absUploads, filename), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}
	defer file.Close()

	var (
		iv              = []byte{}
		requirementsMet = len(password) >= 6 &&
			len(password) < 128 &&
			formFile.Size < int64(s.config.MaxEncryptableFileSize)
		confirmation = ctx.PostForm("can_encrypt") == "yes"
	)

	switch {
	case requirementsMet && confirmation:
		encryptionKey, err := cryptography.DeriveBasicKey([]byte(password), s.encryptionKeySize)
		if err != nil {
			api.SetErrResponse(ctx, err)
			return
		}

		iv, err = cryptography.EncryptStream(file, formFileHandle, encryptionKey)
		if err != nil {
			api.SetErrResponse(ctx, err)
			return
		}

	case !requirementsMet && confirmation:
		api.SetErrResponse(ctx, api.ErrInvalidBody)
		return

	default:
		if _, err := io.Copy(file, formFileHandle); err != nil {
			api.SetErrResponse(ctx, err)
			return
		}
	}

	hash, err := s.passwordhash([]byte(password))
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	title := ctx.Request.FormValue("title")
	err = s.db.NewFile(ctx, db.NewFileParams{
		FileUuid:     filename,
		Title:        sql.NullString{String: title, Valid: title != ""},
		Passwdhash:   sql.NullString{String: string(hash), Valid: string(hash) != ""},
		EncryptionIv: iv,
		AccessToken:  uuid.NewString(),
		FileSize:     int32(formFile.Size),
		Encrypted:    len(iv) > 0,
	})
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	if err = s.linkToAccount(ctx, filename); err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	session, err := s.auth.Register(ctx, filename)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	s.persistSession(ctx, formatauthcookiename(session.filename), session.Token)

	ctx.JSON(http.StatusOK, gin.H{
		"filename": filename,
	})
}
