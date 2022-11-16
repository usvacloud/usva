package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/api/middleware"
	"github.com/romeq/usva/dbengine"
	"github.com/romeq/usva/utils"
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

	a, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		a = 5432
	}

	dbengine.Init(
		uint16(a),
		utils.StringOr(os.Getenv("DB_HOST"), "127.0.0.1"),
		utils.StringOr(os.Getenv("DB_USERNAME_TESTS"), "usva_tests"),
		utils.StringOr(os.Getenv("DB_NAME_TESTS"), "usva_tests"),
		utils.StringOr(os.Getenv("DB_PASSWORD_TESTS"), "testrunner"),
	)

	type payload struct {
		fileData string
		maxSize  int
	}

	tests := []struct {
		name               string
		payload            payload
		expectedCode       int
		verifyResponseJSON bool
	}{
		{
			name: "test1",
			payload: payload{
				fileData: "hello",
				maxSize:  8,
			},
			expectedCode:       200,
			verifyResponseJSON: true,
		},
		{
			name: "test2",
			payload: payload{
				fileData: "hello",
				maxSize:  2,
			},
			expectedCode:       413,
			verifyResponseJSON: false,
		},
	}

	for i, tt := range tests {
		responseStruct := struct {
			Filename string
			Message  string
		}{}

		c, r := prepareMultipartBody(t, tt.payload.fileData)
		handler := UploadFile(&middleware.Ratelimiter{}, &APIConfiguration{
			MaxSingleUploadSize: int64(tt.payload.maxSize),
			UploadsDir:          t.TempDir(),
		})
		handler(c)

		// make sure the test ran correctly
		assert.EqualValues(t, tt.expectedCode, r.Code)

		if tt.verifyResponseJSON {
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
