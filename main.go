// should read: https://rcrowley.org/talks/strange-loop-2013.html#1

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"goweb/dbwrapper"
	"goweb/handlers"
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

	db, err := dbwrapper.Connect("postgres", dsn+"?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Generate Schema
	log.Println("Generate DB Schema...")
	if _, err := db.Exec(schema); err != nil {
		log.Println(err)
	}

	// Create our logger
	logger := log.New(os.Stdout, "", 0)

	hDB := &handlers.DB{DB: db.DB}
	r := mux.NewRouter()

	// my version of 'HTTP closure'
	config := "my config"
	r.HandleFunc("/welcome/{name}", handlers.WelcomeHandler(config))
	r.HandleFunc("/_welcome/{name}", handlers.Middleware(logger, handlers.WelcomeHandler(config)))

	// hanlder with no closure
	r.HandleFunc("/about", handlers.About)

	r.Handle("/users", hDB.UserList()).Methods("GET", "HEAD")
	r.Handle("/user/{name}", hDB.UserHandler())
	r.Handle("/gendata", hDB.GenDataHandler()).Methods("GET")

	r.Handle("/_user/{name}", handlers.WithMetrics(logger, hDB.UserHandler()))

	log.Fatal(http.ListenAndServe(":8080", r))
}
