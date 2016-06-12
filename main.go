// should read: https://rcrowley.org/talks/strange-loop-2013.html#1

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

type Tx struct {
	*sqlx.Tx
}

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

type Person struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

type Message struct {
	Text string `json:"msg"`
}

func main() {
	dsn := fmt.Sprintf("%s://%s:%s@%s/%s",
		"postgres",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	db, err := sqlx.Connect("postgres", dsn+"?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if _, err := db.Exec(schema); err != nil {
		log.Println(err)
	}

	mydb := &DB{db}

	// mydb, err := Connect("postgres", dsn+"?sslmode=disable")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer mydb.Close()

	// Create our logger
	logger := log.New(os.Stdout, "", 0)

	r := mux.NewRouter()

	// my version of 'HTTP closure'
	config := "my config"
	r.HandleFunc("/welcome/{name}", welcomeHandler(config))
	r.HandleFunc("/_welcome/{name}", middleware(logger, welcomeHandler(config)))

	// hanlder with no closure
	r.HandleFunc("/about", about)

	r.Handle("/users", mydb.UserList()).Methods("GET", "HEAD")
	r.Handle("/user/{name}", mydb.UserHandler())
	r.Handle("/gendata", mydb.GenDataHandler()).Methods("GET")

	r.Handle("/_user/{name}", withMetrics(logger, userHandler()))

	http.ListenAndServe(":8080", r)
}

func (db *DB) UserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		tx, err := db.Beginx()
		if err != nil {
			log.Println(err)
		}
		tx.CreatePerson(&Person{FirstName: name})
		tx.Commit()
	})
}

func (db *DB) UserList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var users []Person
		err := db.Select(&users, "SELECT * FROM person")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
	})
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

// Beginx starts an returns a new transaction.
func (db *DB) Beginx() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func (db *DB) GenDataHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := db.Beginx()
		if err != nil {
			log.Println(err)
		}
		tx.GenerateData()
		if err := tx.Commit(); err != nil {
			log.Println(err)
		}
	})
}

func (tx *Tx) GenerateData() {
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// _, err := tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email) RETURNING id", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	// if err != nil {
	// 	log.Println(err)
	// }
}

func (tx *Tx) CreatePerson(p *Person) error {
	// Validate the input
	if p == nil {
		return errors.New("person required")
	}
	_, err := tx.Exec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", p.FirstName, p.LastName, p.Email)
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// _, err := tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", p)
	return err
}

// my version copied from tsenart's, looks like more of a mess but it works!
func welcomeHandler(config string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "Welcome!, config: %s, user: %s", config, name)
	}
}

func middleware(l *log.Logger, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next(w, r)
		l.Printf("%s: %s %s took %s", time.Now(), r.Method, r.URL, time.Since(began))
	}
}

// took from the awesome: https://gist.github.com/tsenart/5fc18c659814c078378d
func userHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "user: %s", name)
	})
}

func withMetrics(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next.ServeHTTP(w, r)
		l.Printf("%s: %s %s took %s", time.Now(), r.Method, r.URL, time.Since(began))
	})
}

func about(w http.ResponseWriter, r *http.Request) {
	m := Message{"go API, build v0.0.001.992."}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Println(err)
	}
}
