package core

import (
	"github.com/ryankurte/authplz/api"
)

// Interface for a user control module
type UserControlInterface interface {
	// Login method, returns api.LoginStatus result, user interface for further use, error in case of failure
	Login(email string, password string) (*api.LoginStatus, interface{}, error)
}

// Interface for token (creation and?) validation
type TokenControlInterface interface {
	ValidateToken(userid string, tokenString string) (*api.TokenAction, error)
}

// Interface for 2 factor authentication modules
// These modules must inform the login handler as to whether
// further authentication is supported
type SecondFactorInterface interface {
	// Check whether a user can use this 2fa module
	// This depends on what second factors are registered
	IsSupported(userid string) bool
}

// Interface for token handler modules
// These modules accept a token action and user id to execute a task
// For example, the user module accepts 'activate' and 'unlock' actions
type TokenHandlerInterface interface {
	HandleToken(u interface{}, tokenAction api.TokenAction) error
}

// Interface for login handler modules
// These modules are called in the login chain and can evaluate (and reject) login
// attempts.
type LoginHandlerInterface interface {
	CheckLogin(userid string, u interface{}) (bool, error)
}

// Interface for event handler modules
// These modules are bound into the event manager to provide asynchronous services
// based on system events.
// For example, the mailer module accepts a variety of user events and sends mail in response.
type EventHandlerInterface interface {
	handleEvent(userid string, u interface{}) error
}

// Interface for user instances
type UserInterface interface {
	GetExtId() string
	GetEmail() string
}
