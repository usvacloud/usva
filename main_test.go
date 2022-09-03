package main

import (
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/romeq/tapsa/config"
	"github.com/romeq/tapsa/dbengine"
	"github.com/stretchr/testify/assert"
)

func initdb(t *testing.T) {
	f, err := os.Create("test.db")
	if err != nil {
		t.Error(err)
	}
	f.Close()
	dbengine.Init("test.db")
}

func TestUpload(t *testing.T) {
	initdb(t)

	cfg := config.New("127.0.0.1", 8080, []string{""}, false, 20, "uploads")
	testfileName := path.Join(cfg.Files.UploadsDir, "testfile")
	pr, pw := io.Pipe()

	// setup writers
	errhandle(t, os.MkdirAll(cfg.Files.UploadsDir, 0755))
	testfile, err := os.Create(testfileName)
	errhandle(t, err)
	defer testfile.Close()

	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()

		part, err := writer.CreateFormFile("file", testfile.Name())
		errhandle(t, err)
		_, err = part.Write([]byte("testtesttest!"))
		errhandle(t, err)
	}()

	request := httptest.NewRequest("POST", "/api/file/upload", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()
	r := setuprouter(cfg)
	r.ServeHTTP(response, request)

	// test output
	t.Log(response.Body)
	assert.Equal(t, 200, response.Code)

	uploadedFile, err := os.Open("uploads/testfile")
	errhandle(t, err)
	fileEquals(t, testfile, uploadedFile)

	t.Cleanup(func() {
		err := os.Remove("test.db")
		errhandle(t, err)
		err = os.Remove("uploads/testfile")
		errhandle(t, err)
	})
}

func fileEquals(t *testing.T, z *os.File, f *os.File) {
	zc, err := io.ReadAll(z)
	errhandle(t, err)
	fc, err := io.ReadAll(f)
	errhandle(t, err)
	assert.Equal(t, zc, fc)
}

func errhandle(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
