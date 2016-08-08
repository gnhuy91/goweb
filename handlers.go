package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gnhuy91/goweb/models"

	"github.com/gorilla/mux"
)

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func UserHandler(db *DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			vars := mux.Vars(r)
			userIDStr := vars["id"]
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				log.Println(err)
				http.Error(w, "invalid user id", http.StatusBadRequest)
				return
			}

			user, err := models.UserInfoByID(db, userID)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(user); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case "POST":
			var u models.User
			u = models.NewUser()

			err := json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if *u.(*models.UserInfo) == (models.UserInfo{}) {
				http.Error(w, "user is empty", http.StatusBadRequest)
				return
			}

			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = u.Insert(tx)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = tx.Commit()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

		case "PUT":
			vars := mux.Vars(r)
			userIDStr := vars["id"]
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				log.Println(err)
				http.Error(w, "invalid user id", http.StatusBadRequest)
				return
			}

			// Parse body
			var u models.UserInfo
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				log.Println(err)
				http.Error(w, "unknown format", http.StatusBadRequest)
				return
			}
			if u == (models.UserInfo{}) {
				http.Error(w, "user is empty", http.StatusBadRequest)
				return
			}

			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			u.ID = userID
			if err := u.Upsert(tx); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := tx.Commit(); err != nil {
				log.Println(err)
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
			}

		case "DELETE":
			vars := mux.Vars(r)
			userIDStr := vars["id"]
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				log.Println(err)
				http.Error(w, "invalid user id", http.StatusBadRequest)
				return
			}

			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			u := models.UserInfo{ID: userID}
			if err := u.Delete(tx); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err := tx.Commit(); err != nil {
				log.Println(err)
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func UserList(db *DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			users, err := models.UserInfoAll(db)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(users); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		case "POST":
			var users models.Users
			err := json.NewDecoder(r.Body).Decode(&users)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// simple data validation
			if users == nil || len(users) == 0 {
				http.Error(w, "user is empty", http.StatusBadRequest)
				return
			}
			for _, u := range users {
				if u == nil || u.Email == "" || u.FirstName == "" || u.LastName == "" {
					w.WriteHeader(http.StatusBadRequest)
				}
			}

			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := users.Insert(tx); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := tx.Commit(); err != nil {
				log.Println(err)
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	})
}

func GenDataHandler(db *DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := db.Begin()
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
	tx.Exec("INSERT INTO user_info (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.Exec("INSERT INTO user_info (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
}

func WithMetrics(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next.ServeHTTP(w, r)
		l.Printf("%s: %s %s took %s", time.Now(), r.Method, r.URL, time.Since(began))
	})
}

func About(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{"text": "go API, build v0.0.001.992."}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Println(err)
	}
}
