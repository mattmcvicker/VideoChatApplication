package handlers

import (
	"github.com/mattmcvicker/WeFeud/servers/gateway/models/users"
	"github.com/mattmcvicker/WeFeud/servers/gateway/sessions"
)

// HandlerContext is a reciever for HTTP handler functions
type HandlerContext struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
}

// NewHandlerContext is a constructor for handler context receiver
func NewHandlerContext(sessionKey string, sessionStore sessions.Store, userStore users.Store) *HandlerContext {
	return &HandlerContext{sessionKey, sessionStore, userStore}
}
