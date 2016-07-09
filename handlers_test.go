package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gnhuy91/goweb/models"
)

var db *DB

func TestMain(m *testing.M) {
	dbc, err := Connect(dbDriver, configDSN())
	if err != nil {
		panic(err)
	}
	// assign to global var so following tests can make use of it
	db = dbc
	defer db.Close()

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() {
	fmt.Println("Create DB Schema...")
	if _, err := db.Exec(schema); err != nil {
		fmt.Println(err)
	}
}

func teardown() {
	fmt.Println("Drop DB Schema...")

	schema := `DROP TABLE user_info`
	_, err := db.Exec(schema)
	if err != nil {
		fmt.Println(err)
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
		// usage of GenerateHandlerTester, not so useful incase
		// we need to modified the request headers.
		tester := GenerateHandlerTester(t, NewRouter(db))
		rec := tester(method, url, body)

		errMsg := "%s %s, body: %s - want %v, got %v"
		errVars := []interface{}{method, url, body, code, rec.Code}

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

		NewRouter(db).ServeHTTP(rec, req)
		errMsg := "%s %s, body: %s - want %v, got %v"
		errVars := []interface{}{method, url, body, code, rec.Code}

		if rec.Code != code {
			t.Errorf(errMsg, errVars...)
		}
	}
}

func TestUserList_StatusOK(t *testing.T) {
	const (
		url    = "/users"
		method = "GET"
		code   = http.StatusOK
	)

	req, _ := http.NewRequest(method, url, nil)
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	errMsg := "%s %s, want %v, got %v"
	errVars := []interface{}{method, url, code, rec.Code}

	if rec.Code != code {
		t.Errorf(errMsg, errVars...)
	}
}

func TestUserUpdate(t *testing.T) {
	const (
		userID = 1
		method = "PUT"
		code   = http.StatusOK
	)
	var url = "/user/" + strconv.Itoa(userID)

	body := `{
			"first_name": "Huy",
			"last_name": "Giang",
			"email": "abc@gmail.com"
		}`

	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	rec := httptest.NewRecorder()

	NewRouter(db).ServeHTTP(rec, req)
	errMsg := "%s %s, body: %s - want %v, got %v"
	errVars := []interface{}{method, url, body, code, rec.Code}

	if rec.Code != code {
		t.Errorf(errMsg, errVars...)
	}

	// Check if the Update took effect by query the user by ID
	// and then compare it with the test body.
	// Both should be parsed to the struct to be able to compare.
	rec.Flush()
	req, _ = http.NewRequest("GET", url, nil)
	NewRouter(db).ServeHTTP(rec, req)

	var userFromTest, userFromDB models.UserInfo
	json.NewDecoder(rec.Body).Decode(&userFromDB)

	// manually assign ID here since PUT get ID from url path, not req body
	userFromTest.ID = userID
	json.NewDecoder(strings.NewReader(body)).Decode(&userFromTest)

	if userFromDB != userFromTest {
		t.Errorf("Update %s went wrong, want %+v, got %+v", url, userFromTest, userFromDB)
	}
}
