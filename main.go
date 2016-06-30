// should read: https://rcrowley.org/talks/strange-loop-2013.html#1

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"goweb/handlers"
)

var schema = `
CREATE TABLE person (
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

	db, err := sqlx.Connect("postgres", dsn+"?sslmode=disable")
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

	mydb := &handlers.DB{DB: db}
	r := mux.NewRouter()

	// my version of 'HTTP closure'
	config := "my config"
	r.HandleFunc("/welcome/{name}", handlers.WelcomeHandler(config))
	r.HandleFunc("/_welcome/{name}", handlers.Middleware(logger, handlers.WelcomeHandler(config)))

	// hanlder with no closure
	r.HandleFunc("/about", handlers.About)

	r.Handle("/users", mydb.UserList()).Methods("GET", "HEAD")
	r.Handle("/user/{name}", mydb.UserHandler())
	r.Handle("/gendata", mydb.GenDataHandler()).Methods("GET")

	r.Handle("/_user/{name}", handlers.WithMetrics(logger, mydb.UserHandler()))

	http.ListenAndServe(":8080", r)
}
