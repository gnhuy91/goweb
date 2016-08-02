package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	svc "github.com/gnhuy91/goweb"
	"github.com/gnhuy91/goweb/models"
)

type testParam struct {
	BaseURL string
	UserID  int
	Method  string
	Code    int
}

func (tp *testParam) ReqURL() string {
	u, _ := svc.ConcatURL(tp.BaseURL, strconv.Itoa(tp.UserID))
	return u
}

func TestUserInsert_ValidBody(t *testing.T) {
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
		tester := svc.GenerateHandlerTester(t, svc.NewRouter(db, logger))
		rec := tester(method, url, body)

		errMsg := "%s %s, body: %s - want %v, got %v"
		errVars := []interface{}{method, url, body, code, rec.Code}

		if rec.Code != code {
			t.Errorf(errMsg, errVars...)
		}
	}
}

func TestUserInsert_InValidBody(t *testing.T) {
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

		svc.NewRouter(db, logger).ServeHTTP(rec, req)
		errMsg := "%s %s, body: %s - want %v, got %v"
		errVars := []interface{}{method, url, body, code, rec.Code}

		if rec.Code != code {
			t.Errorf(errMsg, errVars...)
		}
	}
}

func TestUserByID_StatusOK(t *testing.T) {
	p := testParam{"/user", 1, "GET", http.StatusOK}

	req, _ := http.NewRequest(p.Method, p.ReqURL(), nil)
	rec := httptest.NewRecorder()

	svc.NewRouter(db, logger).ServeHTTP(rec, req)
	errMsg := "%s %s, want %v, got %v"
	errVars := []interface{}{p.Method, p.ReqURL(), p.Code, rec.Code}

	if rec.Code != p.Code {
		t.Errorf(errMsg, errVars...)
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

	svc.NewRouter(db, logger).ServeHTTP(rec, req)
	errMsg := "%s %s, want %v, got %v"
	errVars := []interface{}{method, url, code, rec.Code}

	if rec.Code != code {
		t.Errorf(errMsg, errVars...)
	}
}

func TestUserUpdate(t *testing.T) {
	const (
		baseURL = "/user"
		method  = "PUT"
		code    = http.StatusOK
	)
	body := `{
		"first_name": "Huy",
		"last_name": "Giang",
		"email": "abc@gmail.com"
	}`

	// PUT existed user
	p1 := testParam{baseURL, 1, method, code}
	// PUT non-existed user
	p2 := testParam{baseURL, 9999, method, code}
	ps := []testParam{p1, p2}

	for _, p := range ps {
		req, _ := http.NewRequest(p.Method, p.ReqURL(), strings.NewReader(body))
		rec := httptest.NewRecorder()

		svc.NewRouter(db, logger).ServeHTTP(rec, req)
		errMsg := "%s %s, body: %s - want %v, got %v"
		errVars := []interface{}{p.Method, p.ReqURL(), body, p.Code, rec.Code}

		if rec.Code != code {
			t.Errorf(errMsg, errVars...)
		}

		// Check if the Update took effect by query the user by ID
		// and then compare it with the test body.
		// Both should be parsed to the struct to be able to compare.
		rec.Flush()
		req, _ = http.NewRequest("GET", p.ReqURL(), nil)
		svc.NewRouter(db, logger).ServeHTTP(rec, req)

		var userFromTest, userFromDB models.UserInfo
		json.NewDecoder(rec.Body).Decode(&userFromDB)

		// manually assign ID here since PUT get ID from url path, not req body
		userFromTest.ID = p.UserID
		json.NewDecoder(strings.NewReader(body)).Decode(&userFromTest)

		if userFromDB != userFromTest {
			t.Errorf("Update %s went wrong, want %+v, got %+v", p.ReqURL(), userFromTest, userFromDB)
		}
	}
}

func TestUserDelete_StatusOK(t *testing.T) {
	const (
		url    = "/user/1"
		method = "DELETE"
		code   = http.StatusOK
	)

	req, _ := http.NewRequest(method, url, nil)
	rec := httptest.NewRecorder()

	svc.NewRouter(db, logger).ServeHTTP(rec, req)
	errMsg := "%s %s, want %v, got %v"
	errVars := []interface{}{method, url, code, rec.Code}

	if rec.Code != code {
		t.Errorf(errMsg, errVars...)
	}
}

func TestUserDelete_ShouldNotExist(t *testing.T) {
	const (
		url    = "/user/1"
		method = "GET"
		code   = http.StatusNotFound
	)

	req, _ := http.NewRequest(method, url, nil)
	rec := httptest.NewRecorder()

	svc.NewRouter(db, logger).ServeHTTP(rec, req)
	errMsg := "%s %s, want %v, got %v"
	errVars := []interface{}{method, url, code, rec.Code}

	if rec.Code != code {
		t.Errorf(errMsg, errVars...)
	}
}

func TestUsersInsert(t *testing.T) {
	// TODO: split this test into 2 tests,
	// one for StatusOK assertion and one for verifying result.
	const (
		url    = "/users"
		method = "POST"
		code   = http.StatusOK
		body   = `[
			{
				"first_name": "Huy",
				"last_name": "Giang",
				"email": "abc@mail.com"
			},
			{
				"first_name": "John",
				"last_name": "Doe",
				"email": "johnd@mail.com"
			}
		]`
	)

	// TRUNCATE is faster than DELETE: https://www.postgresql.org/docs/9.1/static/sql-truncate.html
	// RESTART IDENTITY reset auto-increment sequences (id column).
	// if only TRUNCATE table without RESTART IDENTITY,
	// auto-increment fields will keep their sequences.
	db.MustExec("TRUNCATE TABLE user_info RESTART IDENTITY")

	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	rec := httptest.NewRecorder()

	svc.NewRouter(db, logger).ServeHTTP(rec, req)
	errMsg := "%s %s, body: %s - want %v, got %v"
	errVars := []interface{}{method, url, body, code, rec.Code}

	if rec.Code != code {
		t.Errorf(errMsg, errVars...)
	}

	// Check if the operation took effect by querying /users endpoint
	// and then compare it with the test body.
	// Both will be parsed to structs to be able to compare.
	rec.Flush()
	req, _ = http.NewRequest("GET", url, nil)
	svc.NewRouter(db, logger).ServeHTTP(rec, req)

	var usersFromTest, usersFromDB models.Users
	json.NewDecoder(rec.Body).Decode(&usersFromDB)
	json.NewDecoder(strings.NewReader(body)).Decode(&usersFromTest)
	// manually assign ID here since ID is an auto-increment column
	// and 'id' field should not be included in POST body.
	// Thanks to RESTART IDENTITY to make id count predictable as below.
	for i, u := range usersFromTest {
		u.ID = i + 1 // i starts from 0 but DB id starts from 1
	}

	// compare 2 slices, see http://stackoverflow.com/a/15312182/4328963
	if !reflect.DeepEqual(usersFromTest, usersFromDB) {
		t.Errorf("%s %s failed, want %+v, got %+v", method, url, usersFromTest, usersFromDB)
	}
}
