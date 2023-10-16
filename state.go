package dggchat

import (
	"strings"
	"sync"
)

type state struct {
	sync.RWMutex
	users []User
}

func (s *state) removeUser(nick string) {
	s.Lock()
	defer s.Unlock()

	for i, user := range s.users {
		if strings.EqualFold(user.Nick, nick) {
			s.users = append(s.users[:i], s.users[i+1:]...)
			break
		}
	}
}

func (s *state) addUser(user User) {
	s.Lock()
	defer s.Unlock()

	// If you are not in chat (0 instances of your user), and join,
	// chat backend includes your name in the NAMES command, and ALSO
	// sends a JOIN command with your name. This makes sure we do not
	// include ourself 2 times. Otherwise this check would not be needed.
	for _, u := range s.users {
		if strings.EqualFold(user.Nick, u.Nick) {
			return
		}
	}

	s.users = append(s.users, user)
}

func (s *state) updateUser(user User) {
	s.Lock()
	defer s.Unlock()

	for i, u := range s.users {
		if user.ID == u.ID {
			s.users[i] = user
		}
	}
}

func newState() *state {
	s := &state{
		users: make([]User, 0),
	}
	return s
}
