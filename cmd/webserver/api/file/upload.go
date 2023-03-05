package file

import (
	"database/sql"
	"encoding/base64"
	"io"
	"mime/multipart"
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

type uploadForm struct {
	File       *multipart.FileHeader `form:"file"`
	Title      string                `form:"title"`
	Password   string                `form:"password"`
	CanEncrypt string                `form:"can_encrypt"`
	FileHash   string                `form:"hash"`
}

// UploadFile
func (s *Handler) UploadFile(ctx *gin.Context) {
	body, err := api.BindBodyToStruct[uploadForm](ctx)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	filename := uuid.NewString() + path.Ext(body.File.Filename)

	if s.config.MaxSingleUploadSize > 0 && uint64(body.File.Size) > s.config.MaxSingleUploadSize {
		api.SetErrResponse(ctx, api.ErrTooBigBody)
		return
	}

	pwd, err := base64.RawStdEncoding.DecodeString(body.Password)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	absUploadsPath, err := filepath.Abs(s.config.UploadsDir)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	formFileHandle, err := body.File.Open()
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}

	file, err := os.OpenFile(path.Join(absUploadsPath, filename), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		api.SetErrResponse(ctx, err)
		return
	}
	defer file.Close()

	var (
		iv              = []byte{}
		password        = strings.TrimSpace(string(pwd))
		requirementsMet = len(password) >= 6 && len(password) < 128 && body.File.Size < int64(s.config.MaxEncryptableFileSize)
		confirmation    = body.CanEncrypt == "yes"
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

	if err = s.db.NewFile(ctx, db.NewFileParams{
		FileUuid:     filename,
		Title:        sql.NullString{String: body.Title, Valid: body.Title != ""},
		Passwdhash:   sql.NullString{String: string(hash), Valid: string(hash) != ""},
		EncryptionIv: iv,
		AccessToken:  uuid.NewString(),
		FileSize:     int32(body.File.Size),
		Encrypted:    len(iv) > 0,
		FileHash:     body.FileHash,
	}); err != nil {
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
