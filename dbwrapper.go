package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

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

// RetryConnect perform retrying db Connect if failed by 'retryCount' times,
// wait 1 second between each retry.
func RetryConnect(dataDriver, dataSourceName string, retryCount int) (*DB, error) {
	db, err := Connect(dataDriver, dataSourceName)
	for retryCount > 0 {
		if err != nil {
			retryCount--
			time.Sleep(time.Second)
			return RetryConnect(dataDriver, dataSourceName, retryCount)
		}
		break
	}
	return db, err
}
