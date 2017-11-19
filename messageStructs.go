package dggchat

import (
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
		Sender    string
		Timestamp time.Time
		Message   string
		Features  Features
	}

	message struct {
		Nick      string   `json:"nick"`
		Features  []string `json:"features"`
		Timestamp int64    `json:"timestamp"`
		Data      string   `json:"data"`
	}

	namesMessage struct {
		Connections int    `json:"connectioncount"`
		Users       []User `json:"users"`
	}

	// User represents a user with a list of features
	User struct {
		Nick     string   `json:"nick"`
		Features []string `json:"features"`
	}

	roomAction struct {
		Nick      string   `json:"nick"`
		Features  []string `json:"features"`
		Timestamp int64    `json:"timestamp"`
	}

	// Features contains a list of different user features
	Features []string
)

// HasFeature returns true if a feature is in the features list
func (f Features) HasFeature(s string) bool {
	for _, feature := range f {
		if feature == s {
			return true
		}
	}

	return false
}
