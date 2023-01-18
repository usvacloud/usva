package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/internal/dbengine"
	"github.com/romeq/usva/internal/utils"
	"github.com/stretchr/testify/assert"
)

func prepareMultipartBody(t *testing.T, text string) (*gin.Context, *httptest.ResponseRecorder) {
	requestBody := new(bytes.Buffer)
	mw := multipart.NewWriter(requestBody)

	bodyFile, err := mw.CreateFormFile("file", text)
	if assert.NoError(t, err) {
		_, err = bodyFile.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()

	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest("POST", "/", requestBody)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	return c, r
}

func TestUploadFile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	a, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		a = 5432
	}
	db := dbengine.Init(dbengine.DbConfig{
		Host:        utils.StringOr(os.Getenv("DB_HOST"), "127.0.0.1"),
		Port:        a,
		User:        utils.StringOr(os.Getenv("DB_USERNAME_TESTS"), "usva_tests"),
		Password:    utils.StringOr(os.Getenv("DB_PASSWORD_TESTS"), "testrunner"),
		Name:        utils.StringOr(os.Getenv("DB_NAME_TESTS"), "usva_tests"),
		SslDisabled: true,
	})

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
			name: "test1",
			payload: payload{
				fileData: "hello",
				maxSize:  8,
			},
			expectedCode:       200,
			context:            context.Background(),
			verifyResponseJSON: true,
		},
		{
			name: "test2",
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
		path, err := os.MkdirTemp(os.TempDir(), "usva-tmp")
		if err != nil {
			t.Fatal(err)
		}

		server := NewServer(nil, db, &Configuration{
			UseSecureCookie:     false,
			UploadsDir:          path,
			MaxSingleUploadSize: uint64(tt.payload.maxSize),
		}, 16)
		server.UploadFile(c)

		// make sure the test ran correctly
		assert.EqualValues(t, tt.expectedCode, r.Code)

		if tt.verifyResponseJSON {
			e := json.Unmarshal(r.Body.Bytes(), &responseStruct)
			if e != nil {
				t.Fatal(fmt.Sprintf("test %d failed:", i), e)
			}

			_, e = server.db.FileInformation(tt.context, responseStruct.Filename)
			if e != nil {
				t.Fatal(fmt.Sprintf("test %d failed:", i), e)
			}
		}
	}
}
