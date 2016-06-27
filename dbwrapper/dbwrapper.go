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

// // Beginx starts and returns a new transaction
// func (db *DB) Beginx() (*Tx, error) {
// 	tx, err := db.DB.Beginx()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Tx{tx}, nil
// }
