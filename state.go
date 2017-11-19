package dggchat

import (
	"strings"
	"sync"
)

type state struct {
	sync.RWMutex
	connections int
	users       []User
}

func (s *state) removeUser(nick string) {
	s.Lock()
	defer s.Unlock()

	for i, user := range s.users {
		if strings.EqualFold(user.Nick, nick) {
			s.users = append(s.users[:i], s.users[i+1:]...)
			s.connections--
		}
	}
}

func (s *state) addUser(user User) {
	s.Lock()
	defer s.Unlock()

	s.users = append(s.users, user)
}

func newState() *state {
	s := &state{
		connections: 0,
		users:       make([]User, 0),
	}

	return s
}
