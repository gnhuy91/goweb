package dbwrapper

import "github.com/jmoiron/sqlx"

// DB do something
type DB struct {
	*sqlx.DB
}

// Tx represent something
type Tx struct {
	*sqlx.Tx
}
