package main

import "github.com/gorilla/mux"

// NewRouter initialize router with defined routes in routes.go
func NewRouter(db *DB) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes(db) {
		router.
			Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
