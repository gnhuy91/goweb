package transactions

import (
	"errors"
	"goweb/dbwrapper"
	"goweb/models"
)

type Tx dbwrapper.Tx

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
func (tx *Tx) CreateUser(p *models.User) error {
	// Validate the input
	if p == nil {
		return errors.New("user required")
	}
	_, err := tx.Exec("INSERT INTO user_info (first_name, last_name, email) VALUES ($1, $2, $3)", p.FirstName, p.LastName, p.Email)
	return err
}
