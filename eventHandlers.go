package dggchat

type handlers struct {
	msgHandler  func(Message)
	errHandler  func(string)
	joinHandler func(RoomAction)
	quitHandler func(RoomAction)
	pmHandler   func(PrivateMessage)
}

// AddMessageHandler adds a function that will be called every time a message is received
func (s *Session) AddMessageHandler(fn func(Message)) {
	s.handlers.msgHandler = fn
}

// AddErrorHandler adds a function that will be called every time an error message is received
func (s *Session) AddErrorHandler(fn func(string)) {
	s.handlers.errHandler = fn
}

// AddJoinHandler adds a function that will be called every time a user join the chat
func (s *Session) AddJoinHandler(fn func(RoomAction)) {
	s.handlers.joinHandler = fn
}

// AddQuitHandler adds a function that will be called every time a user quits the chat
func (s *Session) AddQuitHandler(fn func(RoomAction)) {
	s.handlers.quitHandler = fn
}

// AddPMHandler adds a function that will be called every time a private message is received
func (s *Session) AddPMHandler(fn func(PrivateMessage)) {
	s.handlers.pmHandler = fn
}
