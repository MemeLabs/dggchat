package dggchat

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A Session represents a connection to destinygg chat
type Session struct {
	sync.RWMutex

	readOnly bool

	// If true, attempt to reconnect on error
	AttempToReconnect bool

	loginKey string

	listening chan bool

	ws *websocket.Conn
}

// ErrAlreadyOpen is thrown when attempting to open a web socket connection
// on a websocket that is already open
var ErrAlreadyOpen = errors.New("web socket is already open")

var wsUrl = url.URL{Scheme: "wss", Host: "www.destiny.gg", Path: "/ws"}

// Open opens a websocket connection to destinygg chat
func (s *Session) Open() error {
	s.Lock()

	defer s.Unlock()

	if s.ws != nil {
		return ErrAlreadyOpen
	}

	header := http.Header{}

	if !s.readOnly {
		header.Add("Cookie", fmt.Sprintf("authtoken=%s", s.loginKey))
	}

	ws, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), header)
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
	s.Lock()
	defer s.Unlock()

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
