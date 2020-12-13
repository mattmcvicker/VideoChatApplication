package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mattmcvicker/WeFeud/servers/gateway/sessions"
)

// checks for cors headers. Returns an error if any are not present
func CheckHeaders(h *httptest.ResponseRecorder) error {
	acao := h.Header().Get("Access-Control-Allow-Origin")
	if acao != "*" {
		return errors.New("Access-Control-Allow-Origin Header not set")
	}
	acam := h.Header().Get("Access-Control-Allow-Methods")
	if acam != "GET, PUT, POST, PATCH, DELETE" {
		return errors.New("Access-Control-Allow-Methods Header not set")
	}
	acah := h.Header().Get("Access-Control-Allow-Headers")
	if acah != "Content-Type, Authorization" {
		return errors.New("Access-Control-Allow-Headers Header not set")
	}
	aceh := h.Header().Get("Access-Control-Expose-Headers")
	if aceh != "Authorization" {
		return errors.New("Access-Control-Expose-Headers Header not set")
	}
	acma := h.Header().Get("Access-Control-Max-Age")
	if acma != "600" {
		return errors.New("Access-Control-Max-Age Header not set")
	}

	return nil
}
func TestNewCORSUsersHandler(t *testing.T) {
	url := "/v1/users"
	resp := httptest.NewRecorder()
	// create memstore for session store
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	// create dummy user store
	userStore := &TestMySQLStore{}
	req := httptest.NewRequest("GET", url, nil)
	// req.Header.Set(c.headertype, c.headervalue)
	context := NewHandlerContext("test string", sessionStore, userStore)
	handler := http.HandlerFunc(context.UsersHandler)
	cors := NewCORS(handler)

	cors.ServeHTTP(resp, req)
	// check headers
	err := CheckHeaders(resp)
	if err != nil {
		// report error
		t.Errorf("Error when checking headers. Check that all correct headers are present")
	}
}

func TestNewCORSSpecificUserHandler(t *testing.T) {
	url := "/v1/users/1"
	resp := httptest.NewRecorder()
	// create memstore for session store
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	// create dummy user store
	userStore := &TestMySQLStore{}
	req := httptest.NewRequest("GET", url, nil)
	// req.Header.Set(c.headertype, c.headervalue)
	context := NewHandlerContext("test string", sessionStore, userStore)
	handler := http.HandlerFunc(context.SpecificUserHandler)
	cors := NewCORS(handler)

	cors.ServeHTTP(resp, req)
	// check headers
	err := CheckHeaders(resp)
	if err != nil {
		// report error
		t.Errorf("Error when checking headers. Check that all correct headers are present")
	}
}

func TestNewCORSSessionsHandler(t *testing.T) {
	url := "/v1/sessions"
	resp := httptest.NewRecorder()
	// create memstore for session store
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	// create dummy user store
	userStore := &TestMySQLStore{}
	req := httptest.NewRequest("GET", url, nil)
	// req.Header.Set(c.headertype, c.headervalue)
	context := NewHandlerContext("test string", sessionStore, userStore)
	handler := http.HandlerFunc(context.SessionsHandler)
	cors := NewCORS(handler)

	cors.ServeHTTP(resp, req)
	// check headers
	err := CheckHeaders(resp)
	if err != nil {
		// report error
		t.Errorf("Error when checking headers. Check that all correct headers are present")
	}
}

func TestNewCORSSpecificSessionHandler(t *testing.T) {
	url := "/v1/session/mine"
	resp := httptest.NewRecorder()
	// create memstore for session store
	sessionStore := sessions.NewMemStore(time.Hour, time.Minute)
	// create dummy user store
	userStore := &TestMySQLStore{}
	req := httptest.NewRequest("GET", url, nil)
	// req.Header.Set(c.headertype, c.headervalue)
	context := NewHandlerContext("test string", sessionStore, userStore)
	handler := http.HandlerFunc(context.SpecificSessionHandler)
	cors := NewCORS(handler)

	cors.ServeHTTP(resp, req)
	// check headers
	err := CheckHeaders(resp)
	if err != nil {
		// report error
		t.Errorf("Error when checking headers. Check that all correct headers are present")
	}
}
