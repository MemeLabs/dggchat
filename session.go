package dggchat

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// A Session represents a connection to destinygg chat
type Session struct {
	// If true, attempt to reconnect on error
	AttempToReconnect bool

	readOnly  bool
	loginKey  string
	listening chan bool
	ws        *websocket.Conn
	handlers  handlers
	state     *state
}

// ErrAlreadyOpen is thrown when attempting to open a web socket connection
// on a websocket that is already open
var ErrAlreadyOpen = errors.New("web socket is already open")

var wsURL = url.URL{Scheme: "wss", Host: "www.destiny.gg", Path: "/ws"}

// Open opens a websocket connection to destinygg chat
func (s *Session) Open() error {
	if s.ws != nil {
		return ErrAlreadyOpen
	}

	header := http.Header{}

	if !s.readOnly {
		header.Add("Cookie", fmt.Sprintf("authtoken=%s", s.loginKey))
	}

	ws, _, err := websocket.DefaultDialer.Dial(wsURL.String(), header)
	if err != nil {
		return err
	}

	s.ws = ws
	s.listening = make(chan bool)

	go s.listen(s.ws, s.listening)

	return nil
}

// Close cleanly closes the connection and stops running listeners
func (s *Session) Close() error {
	if s.ws == nil {
		return nil
	}

	err := s.ws.Close()
	if err != nil {
		return err
	}

	s.ws = nil

	return nil
}

func (s *Session) listen(ws *websocket.Conn, listening <-chan bool) {
	for {
		_, message, err := s.ws.ReadMessage()
		if err != nil {
			if ws != s.ws {
				return
			}

			err := ws.Close()
			if err != nil {
				return
			}

			s.reconnect()
		}

		mslice := strings.Split(string(message[:]), " ")
		if len(mslice) < 2 {
			continue
		}

		mType := mslice[0]
		mContent := strings.Join(mslice[1:], " ")

		switch mType {
		case "MSG":
			m, err := parseMessage(mContent)
			if err != nil {
				continue
			}

			if s.handlers.msgHandler != nil {
				s.handlers.msgHandler(m)
			}
		case "MUTE":
		case "UNMUTE":
		case "BAN":
		case "UNBAN":
		case "SUBONLY":
		// case "PING":
		// case "PONG":
		case "BROADCAST":
		case "PRIVMSG":
		case "ERR":
			errMessage := strings.Replace(mContent, `"`, "", -1)
			s.handlers.errHandler(errMessage)
		case "NAMES":
			n, err := parseNames(mContent)
			if err != nil {
				continue
			}

			s.state.users = n.Users
			s.state.connections = n.Connections
		case "JOIN":
			ra, err := parseRoomAction(mContent)
			if err != nil {
				continue
			}

			user := User{
				Nick:     ra.Nick,
				Features: ra.Features,
			}

			s.state.addUser(user)
		case "QUIT":
			ra, err := parseRoomAction(mContent)
			if err != nil {
				continue
			}

			s.state.removeUser(ra.Nick)
		}

		select {
		case <-listening:
			return
		default:
		}
	}
}

func (s *Session) reconnect() {
	if !s.AttempToReconnect {
		return
	}

	wait := 1
	for {
		err := s.Open()
		if err == nil || err == ErrAlreadyOpen {
			return
		}

		wait *= 2
		<-time.After(time.Duration(wait) * time.Second)

		if wait > 600 {
			wait = 600
		}
	}
}

// GetUser attempts to find the user in the chat room state
// if the user is found, returns the user and true
// otherwise false is returned as the second parameter
func (s *Session) GetUser(name string) (User, bool) {
	for _, user := range s.state.users {
		if strings.EqualFold(name, user.Nick) {
			return user, true
		}
	}

	return User{}, false
}
