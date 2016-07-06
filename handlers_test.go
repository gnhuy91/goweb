package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var db *DB

func TestMain(m *testing.M) {
	dbc, err := Connect(dbDriver, configDSN())
	if err != nil {
		log.Fatalln(err)
	}
	// assign to global var so following tests can make use of it
	db = dbc
	defer db.Close()

	setup()
	code := m.Run()
	shutdown()

	os.Exit(code)
}

func setup() {
	// prepare things here

	// Generate DB Schema
	log.Println("Generate DB Schema...")
	if _, err := db.Exec(schema); err != nil {
		log.Println(err)
	}
}

func shutdown() {
	// tear-down prepared things here
}

func TestUserList_StatusOK(t *testing.T) {
	url := "/users"

	req, _ := http.NewRequest("GET", url, nil)

	// Use Recorder to record handler's response
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("%s didn't return %v", url, http.StatusOK)
	}
}
