// This file contains helper functions which can be used
// in any projects.
package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// HandlerTester returns an httptest.ResponseRecorder,
// given a method type ("GET", "POST", etc), request url and parameters,
// serve the response against the handler and return the ResponseRecorder.
type HandlerTester func(method, url string, params string) *httptest.ResponseRecorder

// GenerateHandlerTester is a wrapper which returns a HandlerTester func.
func GenerateHandlerTester(t *testing.T, handleFunc http.Handler) HandlerTester {
	return func(method, url string, params string) *httptest.ResponseRecorder {
		req, err := http.NewRequest(
			method,
			url,
			strings.NewReader(params),
		)
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set(
			"Content-Type",
			"application/json",
		)
		req.Body.Close()
		rec := httptest.NewRecorder()
		handleFunc.ServeHTTP(rec, req)
		return rec
	}
}

// From Matt Aimonetti's blog post:
// matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
// Creates a new file upload http request with optional extra params.
func newFileUploadRequest(paramName, path string, params map[string]string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
