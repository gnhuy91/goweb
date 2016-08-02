package main

import (
	"log"

	"github.com/gorilla/mux"
)

// NewRouter initialize router with defined routes in routes.go
func NewRouter(db *DB, logger *log.Logger) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes(db, logger) {
		router.
			Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
