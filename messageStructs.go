package dggchat

import (
	"strings"
	"time"
)

// Constants for different types of features a user can have
const (
	FeatureSubscriber    = "subscriber"
	FeatureBot           = "bot"
	FeatureProtected     = "protected"
	FeatureVIP           = "vip"
	FeatureModerator     = "moderator"
	FeatureAdministrator = "admin"
	FeatureTier2         = "flair1"
	FeatureNotable       = "flair2"
	FeatureTier3         = "flair3"
	FeatureTrusted       = "flair4"
	FeatureContributor   = "flair5"
	FeatureCompChallenge = "flair6"
	FeatureEve           = "flair7"
	FeatureTier4         = "flair8"
	FeatureTwitch        = "flair9"
	FeatureSC2           = "flair10"
	FeatureBot2          = "flair11"
	FeatureBroadcaster   = "flair12"
	FeatureTier1         = "flair13"
	FeatureBirthday      = "flair15"
)

// Constants for different types of errors the chat can return
const (
	ErrorTooManyConnections = "toomanyconnections"
	ErrorProtocol           = "protocolerror"
	ErrorNeedLogin          = "needlogin"
	ErrorNoPermission       = "nopermission"
	ErrorInvalidMessage     = "invalidmsg"
	ErrorMuted              = "muted"
	ErrorSubMode            = "submode"
	ErrorThorttled          = "throttled"
	ErrorDuplicate          = "duplicate"
	ErrorNotFound           = "notfound"
	ErrorNeedBanReason      = "needbanreason"
)

type (
	// Message reprents a normal dgg chat message
	Message struct {
		Sender    User
		Timestamp time.Time
		Message   string
	}

	message struct {
		Nick      string   `json:"nick"`
		Features  []string `json:"features"`
		Timestamp int64    `json:"timestamp"`
		Data      string   `json:"data"`
	}

	// Mute represents (un)mutes issued by chat moderators
	Mute struct {
		Sender    User
		Timestamp time.Time
		Target    User
		// Online indicates whether the target was online when targeted
		Online bool
	}

	// Ban represents (un)bans issued by chat moderators
	Ban struct {
		Sender    User
		Timestamp time.Time
		Target    User
		// Online indicates whether the target was online when targeted
		Online bool
	}

	// Names reprents the initial status message containing user information
	Names struct {
		Connections int    `json:"connectioncount"`
		Users       []User `json:"users"`
	}

	// User represents a user with a list of features
	User struct {
		Nick     string   `json:"nick"`
		Features []string `json:"features"`
	}

	// RoomAction represents a user joining or quitting the chat
	RoomAction struct {
		User      User
		Timestamp time.Time
	}

	roomAction struct {
		Nick      string   `json:"nick"`
		Features  []string `json:"features"`
		Timestamp int64    `json:"timestamp"`
	}

	// PrivateMessage represents a received private message from a user
	PrivateMessage struct {
		User      User
		Message   string
		Timestamp time.Time
		ID        int
	}

	privateMessage struct {
		MessageID int    `json:"messageid"`
		Timestamp int64  `json:"timestamp"`
		Nick      string `json:"nick"`
		Data      string `json:"data"`
	}

	// Broadcast represents a chat broadcast
	Broadcast struct {
		Sender    User
		Timestamp time.Time
		Message   string
	}

	// Ping represents a pong response from the server
	Ping struct {
		Timestamp int64 `json:"timestamp"`
	}

	// SubOnly represents a subonly message from the server
	// if active, only subscribers and some other special user classes are allowed to send messages,
	// until another SubOnly message is received with active set to false
	SubOnly struct {
		Sender    User
		Timestamp time.Time
		Active    bool
	}

	subOnly struct {
		Data      string   `json:"data"`
		Timestamp int64    `json:"timestamp"`
		Nick      string   `json:"nick"`
		Features  []string `json:"features"`
	}
)

// HasFeature returns true if user has given feature
func (u *User) HasFeature(s string) bool {
	for _, feature := range u.Features {
		if feature == s {
			return true
		}
	}

	return false
}

// IsAction returns true if the message was an action (/me)
func (m *Message) IsAction() bool {
	return strings.HasPrefix(m.Message, "/me ")
}
