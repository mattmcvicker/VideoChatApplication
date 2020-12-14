package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	// create new SessionID
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	// save the sessionState to the store
	err = store.Save(sessionID, sessionState)
	if err != nil {
		return InvalidSessionID, err
	}

	// add auth header to the ResponseWriter
	w.Header().Add(headerAuthorization, schemeBearer+sessionID.String())

	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	// get the value of the Authorization header
	h := r.Header.Get(headerAuthorization)

	// get auth query string param if no Authorization header is present
	if h == "" {
		h = r.FormValue(paramAuthorization)
	}

	if h == "" {
		return InvalidSessionID, ErrInvalidScheme
	}

	// check if scheme is valid
	if !strings.Contains(h, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}

	h = strings.Replace(h, schemeBearer, "", 1)

	// validate
	sessID, err := ValidateID(h, signingKey)

	// return SessionID if valid
	if err != nil {
		return InvalidSessionID, err
	}

	return sessID, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	// get the SessionID from the request
	sessID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	// get the data associated with that SessionID from the store
	err = store.Get(sessID, sessionState)
	if err == nil {
		return sessID, nil
	}

	return InvalidSessionID, ErrStateNotFound
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	// get the SessionID from the request
	sessID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	// delete the data associated with it in the store
	err = store.Delete(sessID)
	if err == nil {
		return sessID, nil
	}

	return InvalidSessionID, ErrNoSessionID
}
