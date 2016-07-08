package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestInsertUser_ValidBody(t *testing.T) {
	const (
		url    = "/user"
		method = "POST"
		code   = http.StatusOK
	)

	bodies := []string{`{
			"first_name": "Huy",
			"last_name": "Giang",
			"email": "abc@mail.com"
		}`}

	for _, body := range bodies {
		req, _ := http.NewRequest(method, url, strings.NewReader(body))
		rec := httptest.NewRecorder()

		errMsg := "%s %s, body: %s - want %v, got %v"
		errVars := []interface{}{method, url, body, code, rec.Code}

		NewRouter(db).ServeHTTP(rec, req)
		if rec.Code != code {
			t.Errorf(errMsg, errVars...)
		}
	}
}

func TestInsertUser_InValidBody(t *testing.T) {
	const (
		url    = "/user"
		method = "POST"
		code   = http.StatusBadRequest
	)

	bodies := []string{
		`{}`,
		`{"name": "Huy"}`,
	}

	for _, body := range bodies {
		req, _ := http.NewRequest(method, url, strings.NewReader(body))
		rec := httptest.NewRecorder()

		errMsg := "%s %s, body: %s - want %v, got %v"
		errVars := []interface{}{method, url, body, code, rec.Code}

		NewRouter(db).ServeHTTP(rec, req)
		if rec.Code != code {
			t.Errorf(errMsg, errVars...)
		}
	}
}
