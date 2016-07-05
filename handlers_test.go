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
	db, err := Connect(dbDriver, initDSN())
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Generate Schema
	log.Println("Generate DB Schema...")
	if _, err := db.Exec(schema); err != nil {
		log.Println(err)
	}

	// setup()

	code := m.Run()

	// shutdown()

	os.Exit(code)
}

func setup() {

}

func shutdown() {

}

func TestUserList_StatusOK(t *testing.T) {
	url := "/users"

	// Open our connection and setup our handler
	// db, err := Connect(dbDriver, initDSN())
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer db.Close()

	req, _ := http.NewRequest("GET", url, nil)

	// Use Recorder to record handler's response
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("%s didn't return %v", url, http.StatusOK)
	}
}
