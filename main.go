// should read: https://rcrowley.org/talks/strange-loop-2013.html#1

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

// Create our logger
var logger = log.New(os.Stdout, "", 0)

func main() {
	db, err := Connect(dbDriver, configDSN())
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Generate Schema
	log.Println("Ensure DB Schema created ...")
	allErrors, ok := migrate.UpSync(configDSN(), migrationsDir)
	if !ok {
		log.Println("DB migrate Up failed ...")
		for _, err := range allErrors {
			log.Println(err)
		}
	}

	// r := mux.NewRouter()

	// // my version of 'HTTP closure'
	// config := "my config"
	// r.HandleFunc("/welcome/{name}", WelcomeHandler(config))
	// r.HandleFunc("/_welcome/{name}", Middleware(logger, WelcomeHandler(config)))

	// // hanlder with no closure
	// r.HandleFunc("/about", About)

	// r.Handle("/users", UserList(db)).Methods("GET", "HEAD")
	// r.Handle("/user", UserHandler(db)).Methods("POST")

	// r.Handle("/gendata", GenDataHandler(db)).Methods("GET")

	// log.Fatal(http.ListenAndServe(":8080", r))

	r := NewRouter(db)
	log.Fatal(http.ListenAndServe(configPort(), r))
}
