package models_test

import (
	"database/sql"
	"testing"

	"github.com/gnhuy91/goweb/models"
)

// This struct is used to mock models.DB methods
type DB struct {
	ExecFunc     func(string, ...interface{}) (sql.Result, error)
	QueryFunc    func(string, ...interface{}) (*sql.Rows, error)
	QueryRowFunc func(string, ...interface{}) *sql.Row
	SelectFunc   func(dest interface{}, query string, args ...interface{}) error
	PrepareFunc  func(string) (*sql.Stmt, error)
}

// Pass to struct field function so that we can replace it later,
// this is due to fact that struct fields can be replaced while
// receiver methods cannot.
// Use this trick when you want to change the mocked function's behavior
// on the wire, ie. capturing mocked function's params which was computed
// via previous calls.
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.ExecFunc(query, args)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRowFunc(query, args)
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.QueryFunc(query, args)
}

func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.SelectFunc(dest, query, args)
}

func (db *DB) Prepare(query string) (*sql.Stmt, error) {
	return db.PrepareFunc(query)
}

// Example Test to demonstrate struct methods mocking.
// This test will verify the correctness of Insert's composed SQL query.
func TestUserInsertQuery(t *testing.T) {
	const outp = `INSERT INTO public.user_info (` +
		`first_name, last_name, email` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`
	var q string

	db := &DB{}
	db.QueryRowFunc = func(query string, args ...interface{}) *sql.Row {
		// capture the query
		q = query
		return nil
	}
	// Insert user
	usr := models.UserInfo{}
	usr.Insert(db)

	// verify the query
	if q != outp {
		t.Errorf("wrong query, want %q, got %q", outp, q)
	}
}
