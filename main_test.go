package main

import (
	"encoding/json"
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
		t.Error("Failed to create database file:", err)
	}
	f.Close()
	dbengine.Init("test.db")
}

func TestUpload(t *testing.T) {
	initdb(t)
	setuplogger()

	cfg := config.New("127.0.0.1", 8080, []string{"127.0.0.1"}, false, true, 20, "uploads")
	testFileName := path.Join("tmp", "testfile")
	testFileContent := `testtesttest!`
	pr, pw := io.Pipe()

	// setup writers
	errhandle(t, os.MkdirAll(cfg.Files.UploadsDir, 0755))
	errhandle(t, os.MkdirAll("tmp", 0755))
	testfile, err := os.Create(testFileName)
	errhandle(t, err)
	defer testfile.Close()

	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()

		part, err := writer.CreateFormFile("file", testfile.Name())
		errhandle(t, err)
		_, err = part.Write([]byte(testFileContent))
		errhandle(t, err)
	}()

	r := setuprouter(cfg)
	request := httptest.NewRequest("POST", "/api/file/upload", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	// test output
	assert.Equal(t, 200, response.Code)

	output := map[string]string{}
	errhandle(t, json.Unmarshal(response.Body.Bytes(), &output))
	filepath := path.Join(cfg.Files.UploadsDir, output["filename"])
	uploadedFile, err := os.Open(filepath)
	errhandle(t, err)
	fileContentEquals(t, testFileContent, uploadedFile)

	t.Cleanup(func() {
		err := os.Remove("test.db")
		errhandle(t, err)
		err = os.Remove(path.Join("tmp", "testfile"))
		errhandle(t, err)
		err = os.Remove(filepath)
		errhandle(t, err)
	})
}

func fileContentEquals(t *testing.T, z string, f *os.File) {
	fc, err := io.ReadAll(f)
	errhandle(t, err)
	assert.Equal(t, z, string(fc))
}

func errhandle(t *testing.T, err error) {
	if err != nil {
		t.Fatal("test resulted in an error: ", err)
	}
}
