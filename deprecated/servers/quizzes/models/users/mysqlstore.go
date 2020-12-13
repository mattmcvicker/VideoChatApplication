package users

import (
	"database/sql"
	"time"

	// import sql driver
	_ "github.com/go-sql-driver/mysql"
)

// SQLStore is db used for updating sql database
type SQLStore struct {
	DB *sql.DB
}

// NewSQLStore method
func NewSQLStore(d *sql.DB) *SQLStore {
	return &SQLStore{DB: d}
}

// GetByID returns the User with the given ID
func (ss *SQLStore) GetByID(id int64) (*User, error) {
	db := ss.DB

	rows, err := db.Query("select * from users where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// define variables to create User struct
	var (
		email     string
		passhash  []byte
		username  string
		firstname string
		lastname  string
		photourl  string
	)
	for rows.Next() {
		// read columns in each row into variables
		err := rows.Scan(&id, &email, &passhash, &username, &firstname, &lastname, &photourl)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	u := User{ID: id, Email: email, PassHash: passhash, UserName: username, FirstName: firstname, LastName: lastname, PhotoURL: photourl}
	return &u, nil
}

//GetByEmail returns the User with the given email
func (ss *SQLStore) GetByEmail(email string) (*User, error) {
	db := ss.DB

	rows, err := db.Query("select * from users where email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// define variables to create User struct
	var (
		id        int64
		passhash  []byte
		username  string
		firstname string
		lastname  string
		photourl  string
	)
	for rows.Next() {
		// read columns in each row into variables
		err := rows.Scan(&id, &email, &passhash, &username, &firstname, &lastname, &photourl)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	u := User{ID: id, Email: email, PassHash: passhash, UserName: username, FirstName: firstname, LastName: lastname, PhotoURL: photourl}
	return &u, nil
}

//GetByUserName returns the User with the given Username
func (ss *SQLStore) GetByUserName(username string) (*User, error) {
	db := ss.DB

	rows, err := db.Query("select * from users where username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// define variables to create User struct
	var (
		id        int64
		email     string
		passhash  []byte
		firstname string
		lastname  string
		photourl  string
	)
	for rows.Next() {
		// read columns in each row into variables
		err := rows.Scan(&id, &email, &passhash, &username, &firstname, &lastname, &photourl)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	u := User{ID: id, Email: email, PassHash: passhash, UserName: username, FirstName: firstname, LastName: lastname, PhotoURL: photourl}
	return &u, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ss *SQLStore) Insert(user *User) (*User, error) {
	db := ss.DB
	res, err := db.Exec("INSERT INTO users(email, passhash, username, firstname, lastname, photourl) VALUES(?, ?, ?, ?, ?, ?)", user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)
	if err != nil {
		return nil, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	// return user struct with updated ID field
	user.ID = lastID

	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ss *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	db := ss.DB
	_, err := db.Exec("UPDATE users SET firstname=?, lastname=? WHERE id=?", updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, err
	}
	user, err := ss.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//Delete deletes the user with the given ID
func (ss *SQLStore) Delete(id int64) error {
	db := ss.DB
	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

//Log inserts the user into the database=
func (ss *SQLStore) Log(userid int64, dt time.Time, ip string) error {
	db := ss.DB
	_, err := db.Exec("INSERT INTO logs(userid, signin_time, client_IP) VALUES(?, ?, ?)", userid, dt, ip)
	if err != nil {
		return err
	}

	return nil
}
