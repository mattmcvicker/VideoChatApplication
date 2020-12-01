package handlers

import (
	"time"

	"github.com/mattmcvicker/WeFeud/servers/gateway/models/users"
)

// SessionState struct to keep track of begin time and authenticated user
type SessionState struct {
	Time time.Time   `json:"time"`
	User *users.User `json:"user"`
}
