package main

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
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
	allErrors, ok := migrate.UpSync(configDSN(), migrationsDir)
	if !ok {
		fmt.Println("DB migrate Up failed ...")
		for _, err := range allErrors {
			fmt.Println(err)
		}
	}
}

func teardown() {
	fmt.Println("Drop DB Schema...")
	allErrors, ok := migrate.DownSync(configDSN(), migrationsDir)
	if !ok {
		fmt.Println("DB migrate Down failed ...")
		for _, err := range allErrors {
			fmt.Println(err)
		}
	}
}
