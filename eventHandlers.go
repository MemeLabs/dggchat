package dggchat

type handlers struct {
	msgHandler func(*Message)
	errHandler func(string)
}

// AddMessageHandler adds a function that will be called every time a message is received
func (s *Session) AddMessageHandler(fn func(*Message)) {
	s.handlers.msgHandler = fn
}

// AddErrorHandler adds a function that will be called every time an error message is received
func (s *Session) AddErrorHandler(fn func(string)) {
	s.handlers.errHandler = fn
}
