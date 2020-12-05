package users

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID is a test function for the SQLStore's GetByID
func TestGetByID(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// expected query to be executed
		query := regexp.QuoteMeta("select * from users where id = ?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByEmail is a test function for the SQLStore's GetByEmail
func TestGetByEmail(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet   string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"test@test.com",
			false,
		},
		{
			"User Not Found",
			&User{},
			"test@test.com",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// expected query to be executed
		query := regexp.QuoteMeta("select * from users where email = ?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnError(ErrUserNotFound)

			// Test GetByEmail()
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)

			// Test GetByEmail()
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByUsername is a test function for the SQLStore's GetByUserName
func TestGetByUsername(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name          string
		expectedUser  *User
		usernameToGet string
		expectError   bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"username",
			false,
		},
		{
			"User Not Found",
			&User{},
			"username",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		// expected query to be executed
		query := regexp.QuoteMeta("select * from users where username = ?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.usernameToGet).WillReturnError(ErrUserNotFound)

			// Test GetByUserName()
			user, err := mainSQLStore.GetByUserName(c.usernameToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.usernameToGet).WillReturnRows(row)

			// Test GetByUserName()
			user, err := mainSQLStore.GetByUserName(c.usernameToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestInsert is a test function for the SQLStore's Insert
func TestInsert(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name        string
		user        *User
		expectError bool
	}{
		{
			"User Inserted",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			false,
		},
		{
			"No attributes to insert",
			&User{},
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// expected insert query
		query := regexp.QuoteMeta("INSERT INTO users(email, passhash, username, firstname, lastname, photourl) VALUES(?, ?, ?, ?, ?, ?)")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectExec(query).WithArgs(c.user.Email, c.user.PassHash, c.user.UserName, c.user.FirstName, c.user.LastName, c.user.PhotoURL).WillReturnError(ErrUserNotFound)

			// Test Insert()
			user, err := mainSQLStore.Insert(c.user)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected exec with an expected result
			mock.ExpectExec(query).WithArgs(c.user.Email, c.user.PassHash, c.user.UserName, c.user.FirstName, c.user.LastName, c.user.PhotoURL).WillReturnResult(sqlmock.NewResult(1, 1))

			// Test Insert()
			user, err := mainSQLStore.Insert(c.user)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			// set user ID to the correct one
			if !reflect.DeepEqual(user, c.user) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestUpdate is a test function for the SQLStore's Update
func TestUpdate(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name        string
		user        *User
		update      *Updates
		expectError bool
	}{
		{
			"Update to existing User",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{"first", "last"},
			false,
		},
		{
			"No user to update",
			&User{},
			&Updates{"first", "last"},
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.user.ID,
			c.user.Email,
			c.user.PassHash,
			c.user.UserName,
			c.update.FirstName,
			c.update.LastName,
			c.user.PhotoURL,
		)

		// update query
		query := regexp.QuoteMeta("UPDATE users SET firstname=?, lastname=? WHERE id=?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectExec(query).WithArgs(c.update.FirstName, c.update.LastName, c.user.ID).WillReturnError(ErrUserNotFound)

			// Test Update()
			user, err := mainSQLStore.Update(c.user.ID, c.update)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {

			// Set up an expected exec with the expected result
			mock.ExpectExec(query).WithArgs(c.update.FirstName, c.update.LastName, c.user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

			// set up an expected get query (using GetByID)
			getQuery := regexp.QuoteMeta("select * from users where id = ?")
			mock.ExpectQuery(getQuery).WithArgs(c.user.ID).WillReturnRows(row)

			// Test Update()
			user, err := mainSQLStore.Update(c.user.ID, c.update)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}

			// update first name and last name
			c.user.FirstName = c.update.FirstName
			c.user.LastName = c.update.LastName

			if !reflect.DeepEqual(user, c.user) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestDelete is a test function for the SQLStore's Delete
func TestDelete(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name        string
		user        *User
		idToDelete  int64
		expectError bool
	}{
		{
			"Delete existing user",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"Delete nonexistant user",
			&User{},
			1,
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// expected query to be executed
		query := regexp.QuoteMeta("DELETE FROM users WHERE id=?")

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectExec(query).WithArgs(c.idToDelete).WillReturnError(ErrUserNotFound)

			// Test Delete()
			err := mainSQLStore.Delete(c.idToDelete)
			if err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {

			// Set up an expected exec with the expected result
			mock.ExpectExec(query).WithArgs(c.user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

			// Test Delete()
			err := mainSQLStore.Delete(c.user.ID)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}
