package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mattmcvicker/WeFeud/servers/gateway/models/users"
	"github.com/mattmcvicker/WeFeud/servers/gateway/sessions"
	"golang.org/x/crypto/bcrypt"
)

// UsersHandler handles requests for the "users" resoruce
func (h *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	// respond with http.StatusMehtodNotAllowed (405) error if request is not a POST
	if r.Method != http.MethodPost {
		http.Error(w, "Request invalid: must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	// check that Content-Type starts with 'application/json'
	// if not, respond with status code http.StatusUnsupportedMediaType (415)
	ctype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		http.Error(w, "Invlaid content type: request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}

	// decode request body json into new user struct
	newUser := users.NewUser{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Bad request: json body data is invalid", http.StatusBadRequest)
		return
	}
	// turn new user into User
	u, err := newUser.ToUser()
	if err != nil {
		http.Error(w, "Bad request: json body data is invalid", http.StatusBadRequest)
		return
	}
	// save user in database
	u, err = h.UserStore.Insert(u)
	if err != nil {
		http.Error(w, "Bad Request: Error saving user to database", http.StatusBadRequest)
		return
	}

	// begin a new session for the new user
	sstate := SessionState{time.Now(), u}
	sessions.BeginSession(h.SigningKey, h.SessionStore, sstate, w)

	// add header and status code
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// encode new user profile in the response body
	enc := json.NewEncoder(w)
	if err := enc.Encode(u); err != nil {
		http.Error(w, "Problem encoding JSON", http.StatusInternalServerError)
		return
	}
}

// SpecificUserHandler handles requests for a specific user
func (h *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	// check that the current user is authenticated
	user, err := GetAuthenticatedUser(h, r)
	if err != nil {
		// invalid session; user is not authenticated
		http.Error(w, "User must be authenticated", http.StatusUnauthorized)
		return
	}
	urlID := strings.Split(r.URL.Path, "/")[3]
	if urlID == "me" {
		urlID = strconv.FormatInt(user.ID, 10)
	}
	uid, err := strconv.ParseInt(urlID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid url parameter", http.StatusBadRequest)
		return
	}
	// make sure user exists
	u, err := h.UserStore.GetByID(uid)
	if err != nil {
		http.Error(w, "Invalid Request: No user found with that ID", http.StatusNotFound)
		return
	}

	// GET request
	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// encode new user profile in the response body
		enc := json.NewEncoder(w)
		if err := enc.Encode(u); err != nil {
			http.Error(w, "Problem encoding JSON", http.StatusInternalServerError)
			return
		}

	} else if r.Method == http.MethodPatch {
		// check uid against currently authenticated user id
		if uid != user.ID {
			http.Error(w, "Invalid Request: You are not allowed to view this resource", http.StatusForbidden)
			return
		}
		// check content type header
		ctype := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ctype, "application/json") {
			http.Error(w, "Invlaid content type: request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// decode body into updates struct
		updates := users.Updates{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Error decoding json", http.StatusBadRequest)
			return
		}
		// update the user
		user, err := h.UserStore.Update(uid, &updates)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		// encode user into response body
		enc := json.NewEncoder(w)
		if err := enc.Encode(user); err != nil {
			http.Error(w, "Problem encoding JSON", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

// SessionsHandler handles requests for the "sessions" resource
// also logs when a user tries to sign-in
func (h *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// check content type header
		ctype := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ctype, "application/json") {
			http.Error(w, "Invlaid content type: request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// decode credentials
		credentials := users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Error decoding json", http.StatusBadRequest)
			return
		}
		// find user and authenticate
		// if no user found, take the same amount of time as authenticating
		user, err := h.UserStore.GetByEmail(credentials.Email)
		if user != nil {
			authErr := user.Authenticate(credentials.Password)
			if authErr != nil {
				http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
				return
			}
		} else if err != nil {
			// do something that takes the same amount of time
			bcrypt.CompareHashAndPassword([]byte("password"), []byte("password"))
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}

		// begin new session if authentication is successful
		sstate := SessionState{time.Now(), user}
		sessions.BeginSession(h.SigningKey, h.SessionStore, sstate, w)

		// log user
		err = LogUser(h, r, user)
		if err != nil {
			http.Error(w, "Error logging sign in attempt", http.StatusInternalServerError)
			return
		}

		// respond to client
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")

		// encode user into response body
		enc := json.NewEncoder(w)
		if err := enc.Encode(user); err != nil {
			http.Error(w, "Problem encoding JSON", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificSessionHandler handles requests related to a specific authenticated session
func (h *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		urlsession := strings.Split(r.URL.Path, "/")[3]
		if urlsession != "mine" {
			http.Error(w, "Forbidden request", http.StatusForbidden)
			return
		}
		// end session
		sessions.EndSession(r, h.SigningKey, h.SessionStore)
		// response with plain text message
		w.Write([]byte("signed out\n"))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// GetAuthenticatedUser fetches the session state from the session store and returns authenticated user
func GetAuthenticatedUser(h *HandlerContext, r *http.Request) (*users.User, error) {
	sess := SessionState{}
	_, err := sessions.GetState(r, h.SigningKey, h.SessionStore, &sess)
	if err != nil {
		// invalid session; user is not authenticated
		return nil, err
	}
	return sess.User, nil
}

// LogUser logs all user sing-in attempts to the POST /v1/sessions endpoint
// returns appropriate http Errors
func LogUser(h *HandlerContext, r *http.Request, u *users.User) error {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ip = strings.Split(xff, ", ")[0]
	}

	h.UserStore.Log(u.ID, time.Now(), ip)
	return nil
}

//HandleTestPath handles
func HandleTestPath(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test"))
}
