package users

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestUsersValidate tests user's Validate()
func TestUsersValidate(t *testing.T) {
	users := []NewUser{
		// 0: invalid email
		NewUser{Email: "invalid", Password: "pass123", PasswordConf: "pass123", UserName: "userName", FirstName: "user", LastName: "name"},
		// 1: password too short
		NewUser{Email: "email@uw.edu", Password: "pa", PasswordConf: "pa", UserName: "userName", FirstName: "user", LastName: "name"},
		// 2: password mismatch
		NewUser{Email: "email@uw.edu", Password: "pass123", PasswordConf: "pass456", UserName: "userName", FirstName: "user", LastName: "name"},
		// 3: username empty
		NewUser{Email: "email@uw.edu", Password: "pass123", PasswordConf: "pass123", UserName: "", FirstName: "user", LastName: "name"},
		// 4: username contains space
		NewUser{Email: "email@uw.edu", Password: "pass123", PasswordConf: "pass123", UserName: "user name", FirstName: "user", LastName: "name"},
		// 5: Valid User
		NewUser{Email: "email@uw.edu", Password: "pass123", PasswordConf: "pass123", UserName: "userName", FirstName: "user", LastName: "name"},
	}
	cases := []struct {
		name        string
		hint        string
		User        NewUser
		expectError bool
	}{
		{"Invalid email", "Make sure you are validating email correctly", users[0], true},
		{"Invalid Password", "Make sure you are validating the password length", users[1], true},
		{"Invalid Password", "Make sure you are checking to see if the passwords match", users[2], true},
		{"Invalid Username", "Make sure you are checking that the username is not empty", users[3], true},
		{"Invalid Username", "Make sure you are checking that there are no spaces in the username", users[4], true},
		{"Valid User", "Make sure you are returning a nil error if user is valid", users[5], false},
	}

	for _, c := range cases {
		//test Validating Users
		err := c.User.Validate()
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

// TestUsersToUser tests user's ToUser()
func TestUsersToUser(t *testing.T) {
	invalidNew := NewUser{Email: "invalid", Password: "pass123", PasswordConf: "pass123", UserName: "userName", FirstName: "user", LastName: "name"}
	validNew := NewUser{Email: "email@uw.edu", Password: "pass123", PasswordConf: "pass123", UserName: "userName", FirstName: "user", LastName: "name"}

	cases := []struct {
		name        string
		hint        string
		User        NewUser
		expectError bool
	}{
		{"Invalid User", "Make sure you are validating user properly", invalidNew, true},
		{"Valid User", "Make sure no error is returned if valid user", validNew, false},
	}
	for _, c := range cases {
		//test saving user
		u, err := c.User.ToUser()
		// identical user for checking photoURL when spaces added to email
		c.User.Email = c.User.Email + "  "
		uSpaceAdded, err := c.User.ToUser()
		// identical user for checking photoURL when spaces added to email
		c.User.Email = strings.ToUpper(c.User.Email)
		uAllCaps, err := c.User.ToUser()
		if err != nil && !c.expectError {
			t.Errorf("case %s: unexpected error: %v\nHINT: %s", c.name, err, c.hint)
			// verify that user has correct field names
			if c.User.Email != u.Email {
				t.Errorf("case %s: Email field set incorrectly", c.name)
			}
			if c.User.UserName != u.UserName {
				t.Errorf("case %s: Email field set incorrectly", c.name)
			}
			if c.User.FirstName != u.FirstName {
				t.Errorf("case %s: Email field set incorrectly", c.name)
			}
			if c.User.LastName != u.LastName {
				t.Errorf("case %s: Email field set incorrectly", c.name)
			}
			// check passHash
			if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(c.User.Password)); err != nil {
				t.Errorf("case %s: Password hash was set incorrectly", c.name)
			}
			if uSpaceAdded.PhotoURL != u.PhotoURL {
				t.Errorf("case %s: photoURL was not set correctly", c.name)
			}
			if uAllCaps.PhotoURL != u.PhotoURL {
				t.Errorf("case %s: photoURL was not set correctly", c.name)
			}
		}
		if c.expectError && err == nil {
			t.Errorf("case %s: expected error but didn't get one\nHINT: %s", c.name, c.hint)
		}
	}
}

// TestUsersFullName tests user's FullName()
func TestUsersFullName(t *testing.T) {
	cases := []struct {
		name     string
		hint     string
		User     User
		expected string
	}{
		{"First and Last Name empty", "Make sure function returns an empty string is both names are not set", User{FirstName: "", LastName: ""}, ""},
		{"First name empty", "Make sure that their is no starting space", User{FirstName: "", LastName: "name"}, "name"},
		{"Last name empty", "Make sure that their is no trailing space", User{FirstName: "user", LastName: ""}, "user"},
		{"First and Last name valid", "First and last name should be joined by a space", User{FirstName: "user", LastName: "name"}, "user name"},
	}
	for _, c := range cases {
		name := c.User.FullName()
		if name != c.expected {
			t.Errorf("case %s: \nHINT: %s", c.name, c.hint)
		}
	}
}

// TestUsersSetPasswordAndAuthenticate tests user's Set Password and Authenticate
func TestUsersSetPasswordAndAuthenticate(t *testing.T) {
	u := User{Email: "user@uw.edu", FirstName: "user", LastName: "name"}

	cases := []struct {
		name string
		pass string
	}{
		{"Over 6 character password", "p@ssword"},
		{"Empty password", ""},
	}
	for _, c := range cases {
		err := u.SetPassword(c.pass)

		if err != nil {
			t.Errorf("case: %s: error storing password: %v", c.name, err)
		}

		err = u.Authenticate(c.pass)
		if err != nil {
			t.Errorf("case: %s: Problem Authenticating after password was set: error: %v", c.name, err)
		}

		wrongPass := "password456"
		err = u.Authenticate(wrongPass)
		if err == nil {
			t.Errorf("case: %s. User authenticated with the wrong password expected an error", c.name)
		}
	}
}

// TestUserInvalidPasswordEncryption tests invalid bycrypt password
func TestUsersInvalidPasswordEncryption(t *testing.T) {
	// set invalid cost factor
	bcryptCost = 100
	u := User{Email: "user@uw.edu", FirstName: "user", LastName: "name"}
	pass := ""
	err := u.SetPassword(pass)

	if err == nil {
		t.Errorf("Error expected: password was invalid")
	}
}

// TestUserApplyUpdates tests user's ApplyUpdates
func TestUsersApplyUpdates(t *testing.T) {
	u := User{Email: "user@uw.edu", FirstName: "user", LastName: "name"}
	cases := []struct {
		name     string
		update   Updates
		expected string
	}{
		{"Valid update", Updates{FirstName: "Erin", LastName: "Rochfort"}, "Erin Rochfort"},
		{"Update empty last name", Updates{FirstName: "Erin"}, "user name"},
		{"Update empty first name", Updates{LastName: "Rochfort"}, "user name"},
		{"Update empty strings", Updates{}, "user name"},
	}
	for _, c := range cases {
		u.ApplyUpdates(&c.update)
		name := u.FullName()
		if name != c.expected {
			t.Errorf("case %s: \nExpected: %s but got %s", c.name, c.expected, name)
		}
		// set user back
		u = User{Email: "user@uw.edu", FirstName: "user", LastName: "name"}
	}
}
