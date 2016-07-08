// Package models contains the types for schema 'public'.
package models

import "database/sql"

// DB is the common interface for database operations that can be used with
// types from schema 'public'.
//
// This should work with database/sql.DB and database/sql.Tx.
type DB interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
}

// Log provides the log func used by generated queries.
var Log = func(string, ...interface{}) {}
