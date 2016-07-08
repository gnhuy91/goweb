package main

import "net/http"

type Route struct {
	Name    string
	Methods []string
	Pattern string
	// consider use http.Handler type here so don't have to do type assertion
	// with closure handlers, but in exchange,
	// all handlers func need to be declared in 'closure' func.
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func routes(db *DB) Routes {
	return Routes{
		Route{
			"Index",
			[]string{"GET", "HEAD"},
			"/",
			About,
		},
		Route{
			"About",
			[]string{"GET"},
			"/about",
			About,
		},
		Route{
			"UserByID",
			[]string{"GET", "PUT", "DELETE"},
			"/user/{id}",
			// Because these handlers wrap and return a function,
			// the returned function's return type is unknown at this point,
			// we must do a type assertion so the compiler won't complain
			WithUAA(UserHandler(db)).(http.HandlerFunc),
		},
		Route{
			"User",
			[]string{"POST"},
			"/user",
			UserHandler(db).(http.HandlerFunc),
		},
		Route{
			"UserList",
			[]string{"GET", "POST"},
			"/users",
			WithUAA(UserList(db)).(http.HandlerFunc),
		},
		Route{
			"GenData",
			[]string{"GET"},
			"/gendata",
			WithUAA(GenDataHandler(db)).(http.HandlerFunc),
		},
	}
}
