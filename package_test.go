package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/mattes/migrate/migrate/direction"
	"github.com/mattes/migrate/pipe"
)

var (
	db     *DB
	logger = log.New(os.Stderr, "", 0)
)

func TestMain(m *testing.M) {
	dbc, err := RetryConnect(dbDriver, dsn, dbConnRetryCount)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB after %v attempts - %s", dbConnRetryCount, err))
	}
	fmt.Println("DB connect successful")
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

	p := pipe.New()
	go migrate.Up(p, dsn, migrationsDir)
	ok := writePipe(p)
	if !ok {
		panic("DB migrate Up failed.")
	}
}

func teardown() {
	fmt.Println("\nDrop DB Schema...")

	p := pipe.New()
	go migrate.Down(p, dsn, migrationsDir)
	ok := writePipe(p)
	if !ok {
		panic("DB migrate Down failed.")
	}
}

func writePipe(pipe chan interface{}) (ok bool) {
	okFlag := true
	if pipe != nil {
		for {
			select {
			case item, more := <-pipe:
				if !more {
					return okFlag
				}

				switch item.(type) {
				case string:
					fmt.Println(item.(string))

				case error:
					fmt.Println(item.(error).Error())
					okFlag = false

				case file.File:
					f := item.(file.File)
					if f.Direction == direction.Up {
						fmt.Print(">")
					} else if f.Direction == direction.Down {
						fmt.Print("<")
					}
					fmt.Printf(" %s\n", f.FileName)

				default:
					text := fmt.Sprint(item)
					fmt.Println(text)
				}
			}

			fmt.Println()
		}
	}
	return okFlag
}
