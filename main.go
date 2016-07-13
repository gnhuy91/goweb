// should read: https://rcrowley.org/talks/strange-loop-2013.html#1

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "github.com/lib/pq"
)

// Create our logger
var logger = log.New(os.Stdout, "", 0)

const schema = `
CREATE TABLE IF NOT EXISTS user_info (
	id BIGSERIAL PRIMARY KEY,
	first_name text,
	last_name text,
	email text
);`

func main() {
	db, err := RetryConnect(dbDriver, dsn, dbConnRetryCount)
	if err != nil {
		log.Fatalf("Failed to connect to DB after %v attempts - %s", dbConnRetryCount, err)
	}
	log.Println("DB connect successful")
	defer db.Close()

	// Generate Schema
	log.Println("Ensure DB Schema created ...")
	_, err = db.Exec(schema)
	if err != nil {
		log.Println(err)
	}

	// Init the router
	r := NewRouter(db)
	log.Fatal(http.ListenAndServe(port, r))
}
