package dggchat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A Session represents a connection to destinygg chat.
type Session struct {
	sync.RWMutex
	// If true, attempt to reconnect on error
	attempToReconnect bool

	readOnly bool
	loginKey string
	wsURL    url.URL
	ws       *websocket.Conn
	handlers handlers
	state    *state
	dialer   *websocket.Dialer
}

type messageOut struct {
	Data string `json:"data"`
}

type privateMessageOut struct {
	Nick string `json:"nick"`
	Data string `json:"data"`
}

type muteOut struct {
	Data     string `json:"data"`
	Duration int64  `data:"duration,omitempty"`
}

type banOut struct {
	Nick        string `json:"nick"`
	Reason      string `json:"reason,omitempty"`
	Duration    int64  `json:"duration,omitempty"`
	Banip       bool   `json:"banip,omitempty"`
	Ispermanent bool   `json:"ispermanent"`
}

type pingOut struct {
	Timestamp int64 `json:"timestamp"`
}

// ErrAlreadyOpen is thrown when attempting to open a web socket connection
// on a websocket that is already open.
var ErrAlreadyOpen = errors.New("web socket is already open")

// ErrReadOnly is thrown when attempting to send messages using a read-only session.
var ErrReadOnly = errors.New("session is read-only")

var wsURL = url.URL{Scheme: "wss", Host: "chat.destiny.gg", Path: "/ws"}

// SetURL changes the url that will be used when connecting to the socket server.
// This should be done before calling *session.Open()
func (s *Session) SetURL(u url.URL) {
	s.Lock()
	defer s.Unlock()
	s.wsURL = u
}

// SetDialer changes the websocket dialer that will be used when connecting to the socket server.
func (s *Session) SetDialer(d websocket.Dialer) {
	s.Lock()
	defer s.Unlock()
	s.dialer = &d
}

// Open opens a websocket connection to destinygg chat.
func (s *Session) Open() error {

	s.Lock()
	defer s.Unlock()

	// Only support a single Open() call.
	if s.ws != nil {
		return ErrAlreadyOpen
	}
	return s.open()
}

// call with locks held
func (s *Session) open() error {

	// Repeatedly calling Open() acts like reconnect.
	// this makes sure any old routines die.
	if s.ws != nil {
		_ = s.ws.Close()
		s.ws = nil
	}

	header := http.Header{}
	header.Add("Origin", "https://destiny.gg")
	if !s.readOnly {
		header.Add("Cookie", fmt.Sprintf("authtoken=%s", s.loginKey))
	}

	ws, _, err := s.dialer.Dial(s.wsURL.String(), header)
	if err != nil {
		return err
	}
	s.ws = ws

	go s.listen()

	return nil
}

// Close cleanly closes the connection and stops running listeners
func (s *Session) Close() error {

	s.Lock()
	defer s.Unlock()

	// Assume if Close() is explicitly called, we do not want reconnection behaviour
	s.attempToReconnect = false

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

func (s *Session) reconnect() {

	wait := 1
	for {
		s.Lock()
		err := s.open()
		s.Unlock()

		if err == nil {
			return
		}

		wait *= 2
		if wait > 32 {
			wait = 32
		}
		time.Sleep(time.Duration(wait) * time.Second)
	}
}

func (s *Session) listen() {
	for {
		_, message, err := s.ws.ReadMessage()
		if err != nil {
			if s.handlers.socketErrorHandler != nil {
				s.handlers.socketErrorHandler(err, s)
			}
			if s.attempToReconnect {
				s.reconnect()
			}
			return
		}

		mslice := strings.SplitN(string(message[:]), " ", 2)
		if len(mslice) != 2 {
			continue
		}

		mType := mslice[0]
		mContent := mslice[1]

		switch mType {

		case "MSG":
			m, err := parseMessage(mContent)
			if s.handlers.msgHandler == nil || err != nil {
				continue
			}
			s.handlers.msgHandler(m, s)

		case "MUTE":
			mute, err := parseMute(mContent, s)
			if s.handlers.muteHandler == nil || err != nil {
				continue
			}
			s.handlers.muteHandler(mute, s)

		case "UNMUTE":
			mute, err := parseMute(mContent, s)
			if s.handlers.muteHandler == nil || err != nil {
				continue
			}
			s.handlers.unmuteHandler(mute, s)

		case "BAN":
			ban, err := parseBan(mContent, s)
			if s.handlers.banHandler == nil || err != nil {
				continue
			}
			s.handlers.banHandler(ban, s)

		case "UNBAN":
			ban, err := parseBan(mContent, s)
			if s.handlers.banHandler == nil || err != nil {
				continue
			}
			s.handlers.unbanHandler(ban, s)

		case "SUBONLY":
			so, err := parseSubOnly(mContent)
			if s.handlers.subOnlyHandler == nil || err != nil {
				continue
			}
			s.handlers.subOnlyHandler(so, s)

		case "BROADCAST":
			b, err := parseBroadcast(mContent)
			if s.handlers.broadcastHandler == nil || err != nil {
				continue
			}
			s.handlers.broadcastHandler(b, s)

		case "PRIVMSG":
			pm, err := parsePrivateMessage(mContent, s)
			if s.handlers.pmHandler == nil || err != nil {
				continue
			}
			s.handlers.pmHandler(pm, s)

		case "PRIVMSGSENT":
			// confirms sending of a PM was successful.
			// If not successful, an ERR message is sent anyways. Ignore this.
		case "PING":
		case "PONG":
			p, err := parsePing(mContent)
			if s.handlers.pingHandler == nil || err != nil {
				continue
			}
			s.handlers.pingHandler(p, s)

		case "ERR":
			if s.handlers.errHandler == nil {
				continue
			}
			errMessage := parseErrorMessage(mContent)
			s.handlers.errHandler(errMessage, s)

		case "NAMES":
			n, err := parseNames(mContent)
			if err != nil {
				continue
			}
			s.state.Lock()
			s.state.users = n.Users
			s.state.Unlock()

			if s.handlers.namesHandler != nil {
				s.handlers.namesHandler(n, s)
			}

		case "JOIN":
			ra, err := parseRoomAction(mContent)
			if err != nil {
				continue
			}

			s.state.addUser(ra.User)

			if s.handlers.joinHandler != nil {
				s.handlers.joinHandler(ra, s)
			}

		case "QUIT":
			ra, err := parseRoomAction(mContent)
			if err != nil {
				continue
			}

			s.state.removeUser(ra.User.Nick)

			if s.handlers.quitHandler != nil {
				s.handlers.quitHandler(ra, s)
			}

		case "REFRESH":
			// This message is received immediately before the server closes the
			// connection because user information was changed, and we need to reinitialize.

			// TODO possibly add an eventhandler here
			s.reconnect()
			return
		}
	}
}

// GetUser attempts to find the user in the chat room state.
// If the user is found, returns the user and true,
// otherwise false is returned as the second parameter.
func (s *Session) GetUser(name string) (User, bool) {
	s.state.RLock()
	defer s.state.RUnlock()

	for _, user := range s.state.users {
		if strings.EqualFold(name, user.Nick) {
			return user, true
		}
	}

	return User{}, false
}

// GetUsers returns a list of users currently online
func (s *Session) GetUsers() []User {
	s.state.RLock()
	defer s.state.RUnlock()
	u := make([]User, len(s.state.users))
	copy(u, s.state.users)
	return u
}

func (s *Session) send(message interface{}, mType string) error {
	if s.readOnly {
		return ErrReadOnly
	}
	m, err := json.Marshal(message)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	// Close() might have been called for some reason, this prevents panicing in those cases
	if s.ws == nil {
		return errors.New("connection not established")
	}
	return s.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s %s", mType, m)))
}

// SendMessage sends the given string as a message to chat.
// Note: a return error of nil does not guarantee successful delivery.
// Monitor for error events to ensure the message was sent with no errors.
func (s *Session) SendMessage(message string) error {
	m := messageOut{Data: message}
	return s.send(m, "MSG")
}

// SendMute mutes the user with the given nick.
// If duration is <= 0, the server uses its built-in default duration
func (s *Session) SendMute(nick string, duration time.Duration) error {
	m := muteOut{Data: nick}
	if duration > 0 {
		m.Duration = int64(duration)
	}
	return s.send(m, "MUTE")
}

// SendUnmute unmutes the user with the given nick.
func (s *Session) SendUnmute(nick string) error {
	m := messageOut{Data: nick}
	return s.send(m, "UNMUTE")
}

// SendBan bans the user with the given nick.
// Bans require a ban reason to be specified.
// If duration is <= 0, the server uses its built-in default duration
func (s *Session) SendBan(nick string, reason string, duration time.Duration, banip bool) error {
	b := banOut{
		Nick:   nick,
		Reason: reason,
		Banip:  banip,
	}
	if duration > 0 {
		b.Duration = int64(duration)
	}
	return s.send(b, "BAN")
}

// SendPermanentBan bans the user with the given nick permanently.
// Bans require a ban reason to be specified.
func (s *Session) SendPermanentBan(nick string, reason string, banip bool) error {
	b := banOut{
		Nick:        nick,
		Reason:      reason,
		Banip:       banip,
		Ispermanent: true,
	}
	return s.send(b, "BAN")
}

// SendUnban unbans the user with the given nick.
// Unbanning also removes mutes.
func (s *Session) SendUnban(nick string) error {
	b := messageOut{Data: nick}
	return s.send(b, "UNBAN")
}

// SendAction calls the SendMessage method but also adds
// "/me" in front of the message to make it a chat action
// same caveat with the returned error value applies.
func (s *Session) SendAction(message string) error {
	return s.SendMessage(fmt.Sprintf("/me %s", message))
}

// SendPrivateMessage sends the given user a private message.
func (s *Session) SendPrivateMessage(nick string, message string) error {
	p := privateMessageOut{
		Nick: nick,
		Data: message,
	}
	return s.send(p, "PRIVMSG")
}

// SendSubOnly modifies the chat subonly mode.
// During subonly mode, only subscribers and some other special user classes are allowed to send messages.
func (s *Session) SendSubOnly(subonly bool) error {
	data := "off"
	if subonly {
		data = "on"
	}
	so := messageOut{Data: data}
	return s.send(so, "SUBONLY")
}

// SendBroadcast sends a broadcast message to chat
func (s *Session) SendBroadcast(message string) error {
	b := messageOut{Data: message}
	return s.send(b, "BROADCAST")
}

// SendPing sends a ping to the server with the current timestamp.
func (s *Session) SendPing() error {
	t := pingOut{Timestamp: timeToUnix(time.Now())}
	return s.send(t, "PING")
}
