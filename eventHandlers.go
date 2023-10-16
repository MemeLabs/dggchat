package dggchat

type handlers struct {
	msgHandler          func(Message, *Session)
	pinHandler          func(Pin, *Session)
	namesHandler        func(Names, *Session)
	muteHandler         func(Mute, *Session)
	unmuteHandler       func(Mute, *Session)
	banHandler          func(Ban, *Session)
	unbanHandler        func(Ban, *Session)
	errHandler          func(string, *Session)
	joinHandler         func(RoomAction, *Session)
	quitHandler         func(RoomAction, *Session)
	userUpdateHandler   func(User, *Session)
	pmHandler           func(PrivateMessage, *Session)
	broadcastHandler    func(Broadcast, *Session)
	subscriptionHandler func(Subscription, *Session)
	donationHandler     func(Donation, *Session)
	pingHandler         func(Ping, *Session)
	subOnlyHandler      func(SubOnly, *Session)

	socketErrorHandler func(error, *Session)
}

// AddMessageHandler adds a function that will be called every time a message is received
func (s *Session) AddMessageHandler(fn func(Message, *Session)) {
	s.handlers.msgHandler = fn
}

// AddPinHandler adds a function that will be called every time a pin message is received
func (s *Session) AddPinHandler(fn func(Pin, *Session)) {
	s.handlers.pinHandler = fn
}

// AddNamesHandler adds a function that will be called every time a names message is received
func (s *Session) AddNamesHandler(fn func(Names, *Session)) {
	s.handlers.namesHandler = fn
}

// AddMuteHandler adds a function that will be called every time a mute message is received
func (s *Session) AddMuteHandler(fn func(Mute, *Session)) {
	s.handlers.muteHandler = fn
}

// AddUnmuteHandler adds a function that will be called every time an unmute message is received
func (s *Session) AddUnmuteHandler(fn func(Mute, *Session)) {
	s.handlers.unmuteHandler = fn
}

// AddBanHandler adds a function that will be called every time a ban message is received
func (s *Session) AddBanHandler(fn func(Ban, *Session)) {
	s.handlers.banHandler = fn
}

// AddUnbanHandler adds a function that will be called every time an unban message is received
func (s *Session) AddUnbanHandler(fn func(Ban, *Session)) {
	s.handlers.unbanHandler = fn
}

// AddErrorHandler adds a function that will be called every time an error message is received
func (s *Session) AddErrorHandler(fn func(string, *Session)) {
	s.handlers.errHandler = fn
}

// AddJoinHandler adds a function that will be called every time a user join the chat
func (s *Session) AddJoinHandler(fn func(RoomAction, *Session)) {
	s.handlers.joinHandler = fn
}

// AddQuitHandler adds a function that will be called every time a user quits the chat
func (s *Session) AddQuitHandler(fn func(RoomAction, *Session)) {
	s.handlers.quitHandler = fn
}

// AddUserUpdateHandler adds a function that will be called every time user's information gets updated
func (s *Session) AddUserUpdateHandler(fn func(User, *Session)) {
	s.handlers.userUpdateHandler = fn
}

// AddPMHandler adds a function that will be called every time a private message is received
func (s *Session) AddPMHandler(fn func(PrivateMessage, *Session)) {
	s.handlers.pmHandler = fn
}

// AddBroadcastHandler adds a function that will be called every time a broadcast is sent to the chat
func (s *Session) AddBroadcastHandler(fn func(Broadcast, *Session)) {
	s.handlers.broadcastHandler = fn
}

// AddSubscriptionHandler adds a function that will be called every time a (regular, gifted, or a mass gift) subscription message is received
func (s *Session) AddSubscriptionHandler(fn func(Subscription, *Session)) {
	s.handlers.subscriptionHandler = fn
}

// AddDonationHandler adds a function that will be called every time a donation message is received
func (s *Session) AddDonationHandler(fn func(Donation, *Session)) {
	s.handlers.donationHandler = fn
}

// AddPingHandler adds a function that will be called when a server responds with a pong
func (s *Session) AddPingHandler(fn func(Ping, *Session)) {
	s.handlers.pingHandler = fn
}

// AddSubOnlyHandler adds a function that will will be called every time a subonly message is received
func (s *Session) AddSubOnlyHandler(fn func(SubOnly, *Session)) {
	s.handlers.subOnlyHandler = fn
}

// AddSocketErrorHandler adds a function that will be called every time a socket error occurs
func (s *Session) AddSocketErrorHandler(fn func(error, *Session)) {
	s.handlers.socketErrorHandler = fn
}
