package handlers

import "net/http"

// CORS struct holds handler that middeware calls
type CORS struct {
	handler http.Handler
}

// ServeHTTP handles the request by adding the appropriate headers and
// passing it to the real handler
func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Add("Access-Control-Expose-Headers", "Authorization")
	w.Header().Add("Access-Control-Max-Age", "600")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// call handler
	c.handler.ServeHTTP(w, r)
}

// NewCORS constructs a new CORSE middleware handler
func NewCORS(handlerToWrap http.Handler) *CORS {
	return &CORS{handlerToWrap}
}
