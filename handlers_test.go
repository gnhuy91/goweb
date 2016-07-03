package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserList_StatusOK(t *testing.T) {
	url := "/users"

	// Open our connection and setup our handler
	db, err := Connect(dbDriver, initDSN())
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	req, _ := http.NewRequest("GET", url, nil)

	// Use Recorder to record handler's response
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("%s didn't return %v", url, http.StatusOK)
	}
}
