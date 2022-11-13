package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/config"
	"github.com/romeq/usva/dbengine"
	"github.com/stretchr/testify/assert"
)

func prepareMultipartBody(t *testing.T, text string) (*gin.Context, *httptest.ResponseRecorder) {
	request_body := new(bytes.Buffer)
	mw := multipart.NewWriter(request_body)

	bodyFile, err := mw.CreateFormFile("file", text)
	if assert.NoError(t, err) {
		_, err = bodyFile.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()

	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest("POST", "/", request_body)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	return c, r
}

func Test_uploadFile(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	dbengine.Init(5432, "127.0.0.1", "usva_tests", "usva_tests", "testrunner")

	type payload struct {
		fileData string
		maxSize  int
	}

	tests := []struct {
		name          string
		payload       payload
		expectedCode  int
		verifySucceed bool
	}{
		{
			name: "test1",
			payload: payload{
				fileData: "hello",
				maxSize:  1,
			},
			expectedCode:  200,
			verifySucceed: true,
		},
		{
			name: "test2",
			payload: payload{
				fileData: "hello",
				maxSize:  -1,
			},
			expectedCode:  200,
			verifySucceed: false,
		},
	}

	for i, tt := range tests {
		responseStruct := struct {
			Filename string
			Message  string
		}{}

		c, r := prepareMultipartBody(t, tt.payload.fileData)
		uploadFile(c, &config.Files{
			MaxSize:    int(tt.payload.maxSize),
			UploadsDir: "../test-uploads/",
		})

		// make sure the test ran correctly
		assert.EqualValues(t, tt.expectedCode, r.Code)

		if tt.verifySucceed {
			e := json.Unmarshal(r.Body.Bytes(), &responseStruct)
			if e != nil {
				t.Fatal(fmt.Sprintf("test %d failed:", i), e)
			}

			_, e = dbengine.GetFile(responseStruct.Filename)
			if e != nil {
				t.Fatal(fmt.Sprintf("test %d failed:", i), e)
			}
		}
	}
}