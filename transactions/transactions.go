package transactions

import (
	"errors"
	"goweb/dbwrapper"
	"goweb/models"
)

type DB dbwrapper.DB
type Tx dbwrapper.Tx

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
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

// CreatePerson create a fucking person
func (tx *Tx) CreatePerson(p *models.Person) error {
	// Validate the input
	if p == nil {
		return errors.New("person required")
	}
	_, err := tx.Exec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", p.FirstName, p.LastName, p.Email)
	return err
}
