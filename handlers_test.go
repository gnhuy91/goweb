package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erikstmartin/go-testdb"
)

func TestUserList_StatusOK(t *testing.T) {
	url := "/users"

	// Open fake db from go-testdb
	db, _ := Connect("testdb", "")
	defer db.Close()

	// Stub the query
	query := "select * from user_info"
	columns := []string{"id", "first_name", "last_name", "email"}
	result := ``
	testdb.StubQuery(query, testdb.RowsFromCSVString(columns, result))

	req, _ := http.NewRequest("GET", url, nil)

	// Use Recorder to record handler's response
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("%s didn't return %v", url, http.StatusOK)
	}
}
