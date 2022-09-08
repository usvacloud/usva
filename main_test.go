package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/romeq/tapsa/config"
	"github.com/romeq/tapsa/dbengine"
	"github.com/stretchr/testify/assert"
)

var (
	r        = &gin.Engine{}
	workdir  = ""
	filename = ""
	cfg      = config.Config{
		Server: config.Server{
			Address:        "127.0.0.1",
			Port:           8080,
			TrustedProxies: []string{"127.0.0.1"},
			DebugMode:      false,
			HideRequests:   true,
			AllowedOrigins: []string{"http://127.0.0.1"},
		},
		Files: config.Files{
			MaxSize:    0,
			UploadsDir: "uploads",
		},
	}
)

func prepareWorkspace() {
	setuplogger()

	tmpdir, err := os.MkdirTemp("", "tapsa")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = os.Chdir(tmpdir)
	if err != nil {
		log.Fatal(err.Error())
	}

	workdir = tmpdir
	if err = os.MkdirAll(cfg.Files.UploadsDir, 0755); err != nil {
		log.Fatal(err.Error())
	}

	initdb()
	r = setuprouter(cfg)
}

func initdb() {
	f, err := os.Create("test.db")
	if err != nil {
		log.Fatalln("Failed to create database file:", err)
	}
	f.Close()
	dbengine.Init("test.db")
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

func TestMain(m *testing.M) {
	prepareWorkspace()
	result := m.Run()

	// cleanup
	if err := os.Remove("test.db"); err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(result)
}

func TestUpload(t *testing.T) {
	pr, pw := io.Pipe()
	testFileName := "testfile"
	testFileContent := `jarkko bought a beer; 2.50emarikka bought a dress; 50e`

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

	request := httptest.NewRequest("POST", "/file/upload", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	// assert output
	assert.Equal(t, 200, response.Code)

	output := map[string]string{}
	errhandle(t, json.Unmarshal(response.Body.Bytes(), &output))
	filename = output["filename"]
	filepath := path.Join(cfg.Files.UploadsDir, filename)

	uploadedFile, err := os.Open(filepath)
	errhandle(t, err)
	fileContentEquals(t, testFileContent, uploadedFile)

	t.Cleanup(func() {
		err = os.Remove(path.Join(workdir, "testfile"))
		errhandle(t, err)
	})
}

func TestGet(t *testing.T) {
	requestpath := fmt.Sprintf("/file?filename=%s", filename)
	req := httptest.NewRequest("GET", requestpath, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.GreaterOrEqual(t, 399, res.Code)
}

func TestDelete(t *testing.T) {
	requestpath := fmt.Sprintf("/file?filename=%s", filename)
	req := httptest.NewRequest("DELETE", requestpath, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.GreaterOrEqual(t, 399, res.Code)
}
