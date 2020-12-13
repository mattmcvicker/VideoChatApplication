package users

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	// validate email field
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("%s is an invalid email", nu.Email)
	}

	// validate password: at least 6 characters and matches confirmation
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password is invalid: must be more than 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Password confirmation does not match")
	}

	// validate username: more than 0 length and cannot have spaces
	if len(nu.UserName) <= 0 || strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("Invalid username: cannot be empty or contain spaces")
	}

	return nil
}

// ToUser converts the NewUser to a User, setting the
// PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	// validate
	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	// create gravatar photo url
	h := md5.New()
	email := strings.ToLower(strings.TrimSpace(nu.Email))
	h.Write([]byte(email))
	photoURL := gravatarBasePhotoURL + hex.EncodeToString(h.Sum(nil))

	// create new *User
	newUser := User{Email: nu.Email, UserName: nu.UserName, FirstName: nu.FirstName, LastName: nu.LastName, PhotoURL: photoURL}
	newUser.SetPassword(nu.Password)

	return &newUser, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	full := u.FirstName + " " + u.LastName
	full = strings.TrimSpace(full)

	return full
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}

	u.PassHash = hash

	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password)); err != nil {
		return err
	}

	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	if len(updates.FirstName) <= 0 && len(updates.LastName) <= 0 {
		return errors.New("Invalid update: fields cannot be empty")
	}
	if len(updates.FirstName) > 0 {
		u.FirstName = updates.FirstName
	}
	if len(updates.LastName) > 0 {
		u.LastName = updates.LastName
	}

	return nil
}
