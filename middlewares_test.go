package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUAA_ValidToken(t *testing.T) {
	if uaaURI == "" {
		t.Skip("UAA URI not provided, skipping test")
	}

	url := "/hello"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+uaaAccessToken)

	// Use Recorder to record handler's response
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("%s didn't return %v, return %v", url, http.StatusOK, rec.Code)
	}
}

func TestUAA_InvalidToken(t *testing.T) {
	if uaaURI == "" {
		t.Skip("UAA URI not provided, skipping test")
	}

	url := "/hello"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+uaaAccessToken+"1")

	// Use Recorder to record handler's response
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("%s didn't return %v, return %v", url, http.StatusOK, rec.Code)
	}
}
