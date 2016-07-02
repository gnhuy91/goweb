// should read: https://rcrowley.org/talks/strange-loop-2013.html#1

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE user_info (
	id BIGSERIAL PRIMARY KEY,
	first_name text,
	last_name text,
	email text
);

CREATE TABLE place (
	id BIGSERIAL PRIMARY KEY,
	country text,
	city text NULL,
	telcode integer
)`

// Create our logger
var logger = log.New(os.Stdout, "", 0)

func main() {
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "127.0.0.1:5432"
	}

	dsn := fmt.Sprintf("%s://%s:%s@%s/%s",
		"postgres",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		postgresHost,
		os.Getenv("POSTGRES_DB"))

	db, err := Connect("postgres", dsn+"?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Generate Schema
	log.Println("Generate DB Schema...")
	if _, err := db.Exec(schema); err != nil {
		log.Println(err)
	}

	hDB := &DB{db.DB}
	// r := mux.NewRouter()

	// // my version of 'HTTP closure'
	// config := "my config"
	// r.HandleFunc("/welcome/{name}", WelcomeHandler(config))
	// r.HandleFunc("/_welcome/{name}", Middleware(logger, WelcomeHandler(config)))

	// // hanlder with no closure
	// r.HandleFunc("/about", About)

	// r.Handle("/users", UserList(hDB)).Methods("GET", "HEAD")
	// r.Handle("/user", UserHandler(hDB)).Methods("POST")

	// r.Handle("/gendata", GenDataHandler(hDB)).Methods("GET")

	// log.Fatal(http.ListenAndServe(":8080", r))

	r := NewRouter(hDB)
	log.Fatal(http.ListenAndServe(":8080", r))
}
