package main

import "github.com/jmoiron/sqlx"

// DB do something
type DB struct {
	*sqlx.DB
}

// Tx represent something
type Tx struct {
	*sqlx.Tx
}

// Open returns a DB reference for a data source.
func Open(dataDriver, dataSourceName string) (*DB, error) {
	db, err := sqlx.Open(dataDriver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Connect returns a DB reference for a data source.
func Connect(dataDriver, dataSourceName string) (*DB, error) {
	db, err := sqlx.Connect(dataDriver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
