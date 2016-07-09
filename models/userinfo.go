// Package models contains the types for schema 'public'.
package models

// UserInfo represents a row from 'public.user_info'.
type UserInfo struct {
	ID        int    `db:"id,omitempty" json:"id,omitempty"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
}

// Insert inserts the UserInfo to the database.
func (ui *UserInfo) Insert(db DB) error {
	var err error

	// sql query
	const sqlstr = `INSERT INTO public.user_info (` +
		`first_name, last_name, email` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`

	// run query
	Log(sqlstr, ui.FirstName, ui.LastName, ui.Email)
	err = db.QueryRow(sqlstr, ui.FirstName, ui.LastName, ui.Email).Scan(&ui.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the UserInfo in the database.
func (ui *UserInfo) Update(db DB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE public.user_info SET (` +
		`first_name, last_name, email` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE id = $4`

	// run query
	Log(sqlstr, ui.FirstName, ui.LastName, ui.Email, ui.ID)
	_, err = db.Exec(sqlstr, ui.FirstName, ui.LastName, ui.Email, ui.ID)
	return err
}

// Upsert performs an upsert for UserInfo.
//
// NOTE: PostgreSQL 9.5+ only
func (ui *UserInfo) Upsert(db DB) error {
	var err error

	// sql query
	const sqlstr = `INSERT INTO public.user_info (` +
		`id, first_name, last_name, email` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, first_name, last_name, email` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.first_name, EXCLUDED.last_name, EXCLUDED.email` +
		`)`

	// run query
	Log(sqlstr, ui.ID, ui.FirstName, ui.LastName, ui.Email)
	_, err = db.Exec(sqlstr, ui.ID, ui.FirstName, ui.LastName, ui.Email)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the UserInfo from the database.
func (ui *UserInfo) Delete(db DB) error {
	var err error

	// sql query
	const sqlstr = `DELETE FROM public.user_info WHERE id = $1`

	// run query
	Log(sqlstr, ui.ID)
	_, err = db.Exec(sqlstr, ui.ID)
	if err != nil {
		return err
	}

	return nil
}

// UserInfoByID retrieves a row from 'public.user_info' as a UserInfo.
//
// Generated from index 'user_info_pkey'.
func UserInfoByID(db DB, id int) (*UserInfo, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, first_name, last_name, email ` +
		`FROM public.user_info ` +
		`WHERE id = $1`

	// run query
	Log(sqlstr, id)
	ui := UserInfo{}

	err = db.QueryRow(sqlstr, id).Scan(&ui.ID, &ui.FirstName, &ui.LastName, &ui.Email)
	if err != nil {
		return nil, err
	}

	return &ui, nil
}

// UserInfoAll retrieves all rows from 'public.user_info' as []UserInfo.
func UserInfoAll(db DB) ([]*UserInfo, error) {
	var users []*UserInfo
	err := db.Select(&users, "SELECT * FROM user_info ORDER BY id")
	return users, err
}
