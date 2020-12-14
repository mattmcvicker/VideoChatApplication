package handlers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mattmcvicker/WeFeud/servers/gateway/models/users"
	"github.com/mattmcvicker/WeFeud/servers/gateway/sessions"
)

// define test variables for users
var testNewUserJSON = `{"Email": "test@uw.edu", "Password": "p@ssword1", "PasswordConf": "p@ssword1", "UserName": "testuser", "FirstName": "test", "LastName": "user"}`
var testCredentials = `{"Email": "test@uw.edu", "Password":"p@ssword1"}`
var newUser = users.NewUser{
	Password:     "p@ssword1",
	PasswordConf: "p@ssword1",
	Email:        "test@uw.edu",
	UserName:     "testuser",
	FirstName:    "test",
	LastName:     "user",
}
var testUser, _ = newUser.ToUser()

// define dummy test store for unit tests
type TestMySQLStore struct {
}

func (s *TestMySQLStore) GetByID(id int64) (*users.User, error) {
	if id == 1 {
		testUser.ID = 1
		return testUser, nil
	}
	return nil, errors.New("only id 1 can be found")
}
func (s *TestMySQLStore) GetByEmail(email string) (*users.User, error) {
	if email == "test@uw.edu" {
		testUser.ID = 1
		return testUser, nil
	}
	return nil, errors.New("only email test@uw.edu can be found")
}
func (s *TestMySQLStore) GetByUserName(username string) (*users.User, error) {
	if username == "testuser" {
		testUser.ID = 1
		return testUser, nil
	}
	return nil, errors.New("only username testuser can be found")
}
func (s *TestMySQLStore) Insert(user *users.User) (*users.User, error) {
	return testUser, nil
}
func (s *TestMySQLStore) Update(id int64, updates *users.Updates) (*users.User, error) {
	return testUser, nil
}
func (s *TestMySQLStore) Delete(id int64) error {
	return nil
}
func (s *TestMySQLStore) Log(userid int64, dt time.Time, ip string) error {
	return nil
}

// tests requests made to the users handler
func TestUsersHandler(t *testing.T) {
	url := "/v1/users"
	cases := []struct {
		name           string
		method         string
		headertype     string
		headervalue    string
		body           string
		expectedStatus int
		expectedError  bool
	}{
		{
			"Valid POST request",
			"POST",
			"Content-type",
			"application/json",
			testNewUserJSON,
			http.StatusCreated,
			false,
		},
		{
			"Invalid POST request - Bad Header",
			"POST",
			"Content-type",
			"text/plain",
			testNewUserJSON,
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"Invalid POST request - Bad Body data",
			"POST",
			"Content-type",
			"application/json",
			"bad body form",
			http.StatusBadRequest,
			true,
		},
		{
			"Invalid POST request - Bad New User",
			"POST",
			"Content-type",
			"application/json",
			`{"Email": "test@uw", "Password": "p", "PasswordConf": "p@", "UserName": "tes", "FirstName": "test", "LastName": "user"}`,
			http.StatusBadRequest,
			true,
		},
		{
			"Invalid Request type - GET request",
			"GET",
			"Content-type",
			"application/json",
			testNewUserJSON,
			http.StatusMethodNotAllowed,
			true,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		req := httptest.NewRequest(c.method, url, bytes.NewBuffer([]byte(c.body)))
		req.Header.Set(c.headertype, c.headervalue)
		// create context
		context := NewHandlerContext("test string", sessionStore, userStore)
		context.UsersHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
		if !c.expectedError {
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		}
	}
}

// tests GET requests made to SpecificUserHandler
func TestGetSpecificUserHandler(t *testing.T) {
	key := "test string"
	baseurl := "/v1/users"
	cases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedError  bool
	}{
		{
			"Valid GET request to specific id",
			baseurl + "/1",
			http.StatusOK,
			false,
		},
		{
			"Valid GET request to me url",
			baseurl + "/me",
			http.StatusOK,
			false,
		},
		{
			"Invalid GET request - ID does not exist",
			baseurl + "/12",
			http.StatusNotFound,
			true,
		},
		{
			"Invalid URL for GET request",
			baseurl + "/badurlvalue",
			http.StatusBadRequest,
			true,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		// start session
		sessions.BeginSession(key, sessionStore, SessionState{Time: time.Now(), User: testUser}, resp)
		token := resp.Header().Get("Authorization")
		req := httptest.NewRequest("GET", c.url, nil)
		req.Header.Add("Authorization", token)
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SpecificUserHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
		if !c.expectedError {
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		}
	}
}

// tests PATCH requests made to SpecificUserHandler
func TestPatchSpecificUserHandler(t *testing.T) {
	key := "test string"
	baseurl := "/v1/users"
	cases := []struct {
		name           string
		url            string
		headertype     string
		headervalue    string
		body           string
		expectedStatus int
		expectedError  bool
	}{
		{
			"Valid PATCH request",
			baseurl + "/1",
			"Content-type",
			"application/json",
			`{"FirstName": "first", "LastName": "last"}`,
			http.StatusOK,
			false,
		},
		{
			"Invalid PATCH request - Bad Header",
			baseurl + "/1",
			"Content-type",
			"text/plain",
			`{"FirstName": "first", "LastName": "last"}`,
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"Invalid PATCH request - Bad Body data",
			baseurl + "/1",
			"Content-type",
			"application/json",
			`{"FakeField: "value"}`,
			http.StatusBadRequest,
			true,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		// start session
		sessions.BeginSession(key, sessionStore, SessionState{Time: time.Now(), User: testUser}, resp)
		token := resp.Header().Get("Authorization")
		req := httptest.NewRequest("PATCH", c.url, bytes.NewBuffer([]byte(c.body)))
		req.Header.Set(c.headertype, c.headervalue)
		req.Header.Add("Authorization", token)
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SpecificUserHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
		if !c.expectedError {
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		}
	}
}

// tests PATCH request made to different user in SpecificUserHandler
func TestPatchForbiddenSpecificUserHandler(t *testing.T) {
	key := "test string"
	baseurl := "/v1/users"
	cases := []struct {
		name           string
		url            string
		headertype     string
		headervalue    string
		body           string
		expectedStatus int
		expectedError  bool
	}{
		{

			"Invalid PATCH request - No Permission",
			baseurl + "/1",
			"Content-type",
			"application/json",
			`{"FakeField":"value"}`,
			http.StatusForbidden,
			true,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		// start session
		user := users.User{ID: 2}
		sessions.BeginSession(key, sessionStore, SessionState{Time: time.Now(), User: &user}, resp)
		token := resp.Header().Get("Authorization")
		req := httptest.NewRequest("PATCH", c.url, bytes.NewBuffer([]byte(c.body)))
		req.Header.Set(c.headertype, c.headervalue)
		req.Header.Add("Authorization", token)
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SpecificUserHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
		if !c.expectedError {
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		}
	}
}

// tests GET request made to SpecificUserHandler when unauthorized
func TestUnauthorizedSpecificUserHandler(t *testing.T) {
	key := "test string"
	baseurl := "/v1/users"
	cases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedError  bool
	}{
		{
			"Valid GET request to specific id",
			baseurl + "/1",
			http.StatusUnauthorized,
			true,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		// start session
		req := httptest.NewRequest("GET", c.url, nil)
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SpecificUserHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
		if !c.expectedError {
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		}
	}
}

// tests invalid DELETE request made to SpecificUserHandler
func TestInvlaidRequestSpecificUserHandler(t *testing.T) {
	key := "test string"
	resp := httptest.NewRecorder()
	// create memstore for session store
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	// create dummy user store
	userStore := &TestMySQLStore{}
	// start session
	sessions.BeginSession(key, sessionStore, SessionState{Time: time.Now(), User: testUser}, resp)
	token := resp.Header().Get("Authorization")
	req := httptest.NewRequest("DELETE", "/v1/users/1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Add("Authorization", token)
	// create context
	context := NewHandlerContext("test string", sessionStore, userStore)
	context.SpecificUserHandler(resp, req)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Errorf("incorrect response status code: expected %d but got %d", http.StatusMethodNotAllowed, resp.Code)
	}
}

// tests requests made to the SessionsHandler
func TestSessionsHandler(t *testing.T) {
	key := "test string"
	url := "/v1/sessions"
	cases := []struct {
		name           string
		method         string
		headertype     string
		headervalue    string
		body           string
		expectedStatus int
		expectedError  bool
	}{
		{
			"Valid POST request",
			"POST",
			"Content-type",
			"application/json",
			testCredentials,
			http.StatusCreated,
			false,
		},
		{
			"Invalid POST request - Bad Header",
			"POST",
			"Content-type",
			"text/plain",
			testCredentials,
			http.StatusUnsupportedMediaType,
			true,
		},
		{
			"Invalid POST request - Bad Body data",
			"POST",
			"Content-type",
			"application/json",
			"invalid body form",
			http.StatusBadRequest,
			true,
		},
		{
			"Invalid POST request - Invalid Credentials Nonexistant email",
			"POST",
			"Content-type",
			"application/json",
			`{"Email": "notuser@uw.edu", "Password": "p"}`,
			http.StatusUnauthorized,
			true,
		},
		{
			"Invalid POST request - Invalid Credentials Incorrect Password",
			"POST",
			"Content-type",
			"application/json",
			`{"Email": "test@uw.edu", "Password": "p"}`,
			http.StatusUnauthorized,
			true,
		},
		{
			"Invalid Request type - GET request",
			"GET",
			"Content-type",
			"application/json",
			testCredentials,
			http.StatusMethodNotAllowed,
			true,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		req := httptest.NewRequest(c.method, url, bytes.NewBuffer([]byte(c.body)))
		req.Header.Set(c.headertype, c.headervalue)
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SessionsHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
		if !c.expectedError {
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		}
	}
}

// tests requests made to the SessionsHandler with X-Forwarded-For header
func TestSessionsHandlerSpecialHeader(t *testing.T) {
	key := "test string"
	url := "/v1/sessions"
	cases := []struct {
		name           string
		method         string
		headertype     string
		headervalue    string
		body           string
		expectedStatus int
		expectedError  bool
	}{
		{
			"Valid POST request",
			"POST",
			"Content-type",
			"application/json",
			testCredentials,
			http.StatusCreated,
			false,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		req := httptest.NewRequest(c.method, url, bytes.NewBuffer([]byte(c.body)))
		req.Header.Set(c.headertype, c.headervalue)
		req.Header.Add("X-Forwarded-For", "127.0.0.1, 190.0.2.1")
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SessionsHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
		if !c.expectedError {
			expectedctype := "application/json"
			ctype := resp.Header().Get("Content-Type")
			if len(ctype) == 0 {
				t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
			} else if !strings.HasPrefix(ctype, expectedctype) {
				t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
			}
		}
	}
}

// tests requests made to the SpecificSessionHandler
func TestInvalidSpecificSessionHandler(t *testing.T) {
	key := "test string"
	baseurl := "/v1/sessions"
	cases := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedError  bool
	}{
		{
			"Invalid GET Request",
			"GET",
			baseurl + "/mine",
			http.StatusMethodNotAllowed,
			true,
		},
		{
			"Invalid DELETE Request - bad url",
			"DELETE",
			baseurl + "/1",
			http.StatusForbidden,
			true,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		// start session
		sessions.BeginSession(key, sessionStore, SessionState{Time: time.Now(), User: testUser}, resp)
		token := resp.Header().Get("Authorization")
		req := httptest.NewRequest(c.method, c.url, nil)
		req.Header.Add("Authorization", token)
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SpecificSessionHandler(resp, req)
		if resp.Code != c.expectedStatus {
			t.Errorf("case %s: incorrect response status code: expected %d but got %d", c.name, c.expectedStatus, resp.Code)
		}
	}
}

// tests valid request to SpecificSessionHandler
func TestValidSpecificSessionHandler(t *testing.T) {
	key := "test string"
	baseurl := "/v1/sessions"
	cases := []struct {
		name          string
		method        string
		url           string
		expectedError bool
	}{
		{
			"Valid GET Request",
			"DELETE",
			baseurl + "/mine",
			false,
		},
	}
	for _, c := range cases {
		resp := httptest.NewRecorder()
		// create memstore for session store
		sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
		// create dummy user store
		userStore := &TestMySQLStore{}
		// start session
		sessions.BeginSession(key, sessionStore, SessionState{Time: time.Now(), User: testUser}, resp)
		token := resp.Header().Get("Authorization")
		req := httptest.NewRequest(c.method, c.url, nil)
		req.Header.Add("Authorization", token)
		// create context
		context := NewHandlerContext(key, sessionStore, userStore)
		context.SpecificSessionHandler(resp, req)

		// check to see that signed out was returned
		b, _ := ioutil.ReadAll(resp.Result().Body)
		s := string(b)
		if !strings.Contains(s, "signed out") {
			t.Errorf("case %s: Response does not contain the words 'sign out'", c.name)
		}
		expectedctype := "text/plain"
		ctype := resp.Header().Get("Content-Type")
		if len(ctype) == 0 {
			t.Errorf("No `Content-Type` header found in the response: must be there start with `%s`", expectedctype)
		} else if !strings.HasPrefix(ctype, expectedctype) {
			t.Errorf("incorrect `Content-Type` header value: expected it to start with `%s` but got `%s`", expectedctype, ctype)
		}
	}
}
