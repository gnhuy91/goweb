package models

import (
	"database/sql"
	"errors"
)

type User struct {
	ID        int    `db:"id,omitempty" json:"id,omitempty"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	// xo fields
	_exists, _deleted bool
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

type Message struct {
	Text string `json:"msg"`
}

// Insert inserts the User to the database.
func (u *User) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO user_info (` +
		`first_name, last_name, email` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`

	// run query and get the returned value
	err = db.QueryRow(sqlstr, u.FirstName, u.LastName, u.Email).Scan(&u.ID)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

type XODB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
}
