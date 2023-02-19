package file

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/cmd/webserver/handlers/auth"
	"github.com/romeq/usva/internal/generated/db"
	"github.com/romeq/usva/pkg/cryptography"
)

var protocol = "http"

// UploadFile
func (s *Handler) UploadFile(ctx *gin.Context) {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	filename := uuid.NewString() + path.Ext(formFile.Filename)

	if s.api.MaxSingleUploadSize > 0 && uint64(formFile.Size) > s.api.MaxSingleUploadSize {
		handlers.SetErrResponse(ctx, handlers.ErrTooBigBody)
		return
	}

	formpwd := strings.TrimSpace(ctx.PostForm("password"))
	pwd, err := base64.RawStdEncoding.DecodeString(formpwd)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}
	password := strings.TrimSpace(string(pwd))

	formFileHandle, err := formFile.Open()
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	absUploads, err := filepath.Abs(s.api.UploadsDir)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	file, err := os.OpenFile(path.Join(absUploads, filename), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	var (
		iv              = []byte{}
		requirementsMet = len(password) >= 6 &&
			len(password) < 128 &&
			formFile.Size < int64(s.api.MaxEncryptableFileSize)
		confirmation = ctx.PostForm("can_encrypt") == "yes"
	)

	switch {
	case requirementsMet && confirmation:
		encryptionKey, err := cryptography.DeriveBasicKey([]byte(password), s.encryptionKeySize)
		if err != nil {
			handlers.SetErrResponse(ctx, err)
			return
		}

		iv, err = cryptography.EncryptStream(file, formFileHandle, encryptionKey)
		if err != nil {
			handlers.SetErrResponse(ctx, err)
			return
		}

	case !requirementsMet && confirmation:
		handlers.SetErrResponse(ctx, handlers.ErrInvalidBody)
		return

	default:
		if _, err := io.Copy(file, formFileHandle); err != nil {
			handlers.SetErrResponse(ctx, err)
			return
		}
	}

	hash, err := auth.BCryptPasswordHash([]byte(password))
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	title := ctx.Request.FormValue("title")
	err = s.db.NewFile(ctx, db.NewFileParams{
		FileUuid:     filename,
		Title:        sql.NullString{String: title, Valid: title != ""},
		Passwdhash:   sql.NullString{String: string(hash), Valid: string(hash) != ""},
		EncryptionIv: iv,
		AccessToken:  uuid.NewString(),
		FileSize: sql.NullInt32{
			Int32: int32(formFile.Size),
			Valid: formFile.Size > 0,
		},
	})
	if err != nil {
		handlers.SetErrResponse(ctx, err)
		return
	}

	if strings.Contains(ctx.Request.Header.Get("User-Agent"), "curl") {
		if s.api.UseSecureCookie {
			protocol = "https"
		}

		ctx.String(http.StatusOK, fmt.Sprintf("%s://%s/file/?filename=%s", protocol, s.api.APIDomain, filename))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"filename": filename,
	})
}
