// Package dggchat provides destiny gg chat binding for Go
package dggchat

import (
	"errors"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// ErrTooManyArgs is thrown when a funcion is called with an unexpeted number of arguments
var ErrTooManyArgs = errors.New("function called with unexcepted amount of arguments")

// New creates a new destinygg session. Accepts either 0 or 1 arguments.
// If no login key is provided, a read-only session is returned
func New(args ...string) (*Session, error) {

	if len(args) > 1 {
		return nil, ErrTooManyArgs
	}
	s := &Session{
		attempToReconnect: true,
		state:             newState(),
		dialer:            websocket.DefaultDialer,
	}
	customHost, customHostExist := os.LookupEnv("CUSTOM_WSHOST")
	if customHostExist {
		s.wsURL = url.URL{Scheme: "wss", Host: customHost, Path: "/wss"}
	}
	customOriginHeader, customOriginHeaderExist := os.LookupEnv("CUSTOM_ORIGINHEADER")
	if customOriginHeaderExist {
		s.originHeader = customOriginHeader
	}
	if len(args) == 1 {
		s.loginKey = args[0]
	} else {
		s.readOnly = true
	}

	return s, nil
}
