package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

			user, err := db.GetUserByID(userID)
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
			err := json.NewDecoder(r.Body).Decode(&u)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if u == (models.User{}) {
				http.Error(w, "user is empty", http.StatusBadRequest)
				return
			}

			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tx.CreateUser(&u)
			tx.Commit()

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
			var u models.User
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				log.Println(err)
				http.Error(w, "unknown format", http.StatusBadRequest)
				return
			}
			if u == (models.User{}) {
				http.Error(w, "user is empty", http.StatusBadRequest)
				return
			}

			tx, err := db.Begin()
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := tx.UpdateUserByID(userID, &u); err != nil {
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
			if err := tx.DeleteUserByID(userID); err != nil {
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
			users, err := db.GetUsers()
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
			var users []*models.User
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
			if err := tx.CreateUsers(users); err != nil {
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
	tx.MustExec("INSERT INTO user_info (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO user_info (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	// Named queries can use structs, so if you have an existing struct (i.e. user := &User{}) that you have populated, you can pass it in as &user
	// _, err := tx.NamedExec("INSERT INTO user_info (first_name, last_name, email) VALUES (:first_name, :last_name, :email) RETURNING id", &User{"Jane", "Citizen", "jane.citzen@example.com"})
	// if err != nil {
	// 	log.Println(err)
	// }
}

// CreateUser create a user in the db
func (tx *Tx) CreateUser(m *models.User) error {
	// Validate the input
	if m == nil {
		return errors.New("user required")
	}
	err := m.Insert(tx)
	return err
}

func (tx *Tx) CreateUsers(m []*models.User) error {
	if m == nil {
		return errors.New("user required")
	}

	// Build multiple values query
	query := "INSERT INTO user_info (first_name, last_name, email) VALUES "
	var (
		vals []interface{}
		i    int
	)
	for _, row := range m {
		query += fmt.Sprintf("($%v, $%v, $%v),", i+1, i+2, i+3)
		i += 3
		vals = append(vals, row.FirstName, row.LastName, row.Email)
	}
	// Remove trailing comma
	query = strings.TrimSuffix(query, ",")

	s, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	_, err = s.Exec(vals...)

	return err
}

func (db *DB) GetUsers() ([]*models.User, error) {
	var users []*models.User
	err := db.Select(&users, "SELECT * FROM user_info ORDER BY id")
	return users, err
}

func (db *DB) GetUserByID(userID int) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM user_info WHERE id=$1", userID)
	return user, err
}

func (tx *Tx) UpdateUserByID(userID int, user *models.User) error {
	_, err := tx.NamedExec(
		`UPDATE user_info
		SET
			first_name=:first_name,
			last_name=:last_name,
			email=:email
		WHERE id=:id`, &models.User{
			ID:        userID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	return err
}

func (tx *Tx) DeleteUserByID(userID int) error {
	// TODO: check if id exist before deleting
	_, err := tx.Exec(`DELETE FROM user_info WHERE id=$1`, userID)
	return err
}

// my version copied from tsenart's, looks like more of a mess but it works!
func WelcomeHandler(config string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "Welcome!, config: %s, user: %s", config, name)
	}
}

func Middleware(l *log.Logger, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		began := time.Now()
		next(w, r)
		l.Printf("%s: %s %s took %s", time.Now(), r.Method, r.URL, time.Since(began))
	}
}

// took from the awesome: https://gist.github.com/tsenart/5fc18c659814c078378d
func MyUserHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "user: %s", name)
	})
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
