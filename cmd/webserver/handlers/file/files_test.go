package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/cmd/webserver/handlers"
	"github.com/romeq/usva/cmd/webserver/handlers/auth"
	"github.com/romeq/usva/internal/dbengine"
	"github.com/romeq/usva/internal/utils"
)

func ensureError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func prepareMultipartBody(t *testing.T, text string) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody := new(bytes.Buffer)
	mw := multipart.NewWriter(requestBody)
	defer mw.Close()

	bodyFile, err := mw.CreateFormFile("file", "file.txt")
	ensureError(t, err)

	_, err = bodyFile.Write([]byte(text))
	ensureError(t, err)

	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest("POST", "/", requestBody)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	return c, r
}

func TestUploadFile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	uploadsPath, err := os.MkdirTemp(os.TempDir(), "usva-tmp")
	if err != nil {
		t.Fatal(err)
	}

	db := dbengine.Init(utils.NewTestDatabaseConfiguration())

	type payload struct {
		fileData string
		maxSize  int
	}

	tests := []struct {
		name               string
		payload            payload
		expectedCode       int
		context            context.Context
		verifyResponseJSON bool
	}{
		{
			name: "test-ok",
			payload: payload{
				fileData: "hello",
				maxSize:  8,
			},
			expectedCode:       200,
			context:            context.Background(),
			verifyResponseJSON: true,
		},
		{
			name: "test-not-ok",
			payload: payload{
				fileData: "hello",
				maxSize:  2,
			},
			expectedCode:       413,
			context:            context.Background(),
			verifyResponseJSON: false,
		},
	}

	for i, tt := range tests {
		responseStruct := struct {
			Filename string
			Message  string
		}{}

		c, r := prepareMultipartBody(t, tt.payload.fileData)

		server := handlers.NewServer(nil, db, &handlers.Configuration{
			UseSecureCookie:     false,
			UploadsDir:          uploadsPath,
			MaxSingleUploadSize: uint64(tt.payload.maxSize),
		}, 16)
		NewFileHandler(server, auth.NewAuthHandler(server)).UploadFile(c)

		if tt.expectedCode != r.Code {
			t.Fatalf("expected %d got %d", tt.expectedCode, r.Code)
		}

		if tt.verifyResponseJSON {
			e := json.Unmarshal(r.Body.Bytes(), &responseStruct)
			if e != nil {
				t.Fatal(fmt.Sprintf("test %d failed:", i), e)
			}

			_, e = server.DB.FileInformation(tt.context, responseStruct.Filename)
			if e != nil {
				t.Fatal(fmt.Sprintf("test %d failed:", i), e)
			}
		}
	}
}
