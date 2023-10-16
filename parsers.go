package dggchat

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

func parseMessage(s string) (Message, error) {
	var m message
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return Message{}, err
	}

	user := User{
		ID:          m.ID,
		Nick:        m.Nick,
		Features:    m.Features,
		CreatedDate: m.CreatedDate,
	}

	message := Message{
		Sender:    user,
		Timestamp: unixToTime(m.Timestamp),
		Message:   m.Data,
	}

	return message, nil
}

func parsePin(s string) (Pin, error) {
	var p pin
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return Pin{}, err
	}

	user := User{
		ID:          p.ID,
		Nick:        p.Nick,
		Features:    p.Features,
		CreatedDate: p.CreatedDate,
	}

	pin := Pin{
		Sender:    user,
		Timestamp: unixToTime(p.Timestamp),
		Message:   p.Data,
		UUID:      p.UUID,
	}

	return pin, nil
}

func parseMute(s string, sess *Session) (Mute, error) {
	m, err := parseMessage(s)
	if err != nil {
		return Mute{}, err
	}

	// Try to get features of target, if they are currently online
	targetNick := m.Message
	u, online := sess.GetUser(targetNick)

	mute := Mute{
		Sender:    m.Sender,
		Timestamp: m.Timestamp,
		Target: User{
			ID:          u.ID,
			Nick:        u.Nick,
			Features:    u.Features,
			CreatedDate: u.CreatedDate,
		},
		Online: online,
	}

	return mute, nil
}

func parseBan(s string, sess *Session) (Ban, error) {
	m, err := parseMessage(s)
	if err != nil {
		return Ban{}, err
	}

	// Try to get features of target, if they are currently online
	targetNick := m.Message
	u, online := sess.GetUser(targetNick)

	ban := Ban{
		Sender:    m.Sender,
		Timestamp: m.Timestamp,
		Target: User{
			ID:          u.ID,
			Nick:        u.Nick,
			Features:    u.Features,
			CreatedDate: u.CreatedDate,
		},
		Online: online,
	}

	return ban, nil
}

func parseNames(s string) (Names, error) {
	var n Names
	err := json.Unmarshal([]byte(s), &n)
	if err != nil {
		return Names{}, err
	}

	return n, nil
}

func parseRoomAction(s string) (RoomAction, error) {
	var ra roomAction

	err := json.Unmarshal([]byte(s), &ra)
	if err != nil {
		return RoomAction{}, err
	}

	roomAction := RoomAction{
		User: User{
			ID:          ra.ID,
			Nick:        ra.Nick,
			Features:    ra.Features,
			CreatedDate: ra.CreatedDate,
		},
		Timestamp: unixToTime(ra.Timestamp),
	}

	return roomAction, nil
}

func parseErrorMessage(s string) string {
	return strings.Replace(s, `"`, "", -1)
}

func parsePrivateMessage(s string, sess *Session) (PrivateMessage, error) {
	var pm privateMessage

	err := json.Unmarshal([]byte(s), &pm)
	if err != nil {
		return PrivateMessage{}, err
	}

	privateMessage := PrivateMessage{
		User: User{
			Nick:     pm.Nick,
			Features: make([]string, 0),
		},
		ID:        pm.MessageID,
		Message:   pm.Data,
		Timestamp: unixToTime(pm.Timestamp),
	}

	u, found := sess.GetUser(privateMessage.User.Nick)
	if found {
		privateMessage.User = u
	}

	return privateMessage, nil
}

func parseBroadcast(s string) (Broadcast, error) {
	var m Broadcast

	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return Broadcast{}, err
	}

	return m, nil
}

func parseSubscription(s string) (Subscription, error) {
	var sub subscription
	if err := json.Unmarshal([]byte(s), &sub); err != nil {
		return Subscription{}, err
	}

	recipient := sub.Recipient
	if recipient.Nick == "" {
		recipient = sub.User
	}

	tier := SubTier{
		Tier:  sub.Tier,
		Label: sub.TierLabel,
	}

	subscription := Subscription{
		Sender:    sub.User,
		Recipient: recipient,
		Timestamp: unixToTime(sub.Timestamp),
		Message:   sub.Data,
		Tier:      tier,
		Quantity:  sub.Quantity,
		UUID:      sub.UUID,
	}

	return subscription, nil
}

func parseDonation(s string) (Donation, error) {
	var dono donation
	if err := json.Unmarshal([]byte(s), &dono); err != nil {
		return Donation{}, err
	}

	donation := Donation{
		Sender:    dono.User,
		Timestamp: unixToTime(dono.Timestamp),
		Message:   dono.Data,
		Amount:    dono.Amount,
		UUID:      dono.UUID,
	}

	return donation, nil
}

func parseSubOnly(s string) (SubOnly, error) {
	var so subOnly

	err := json.Unmarshal([]byte(s), &so)
	if err != nil {
		return SubOnly{}, err
	}

	subonly := SubOnly{
		Sender: User{
			ID:          so.ID,
			Nick:        so.Nick,
			Features:    so.Features,
			CreatedDate: so.CreatedDate,
		},
		Timestamp: unixToTime(so.Timestamp),
		// the backend specifies values "on" and "off" ONLY.
		Active: so.Data == "on",
	}

	return subonly, nil
}

func parsePing(s string) (Ping, error) {
	var p Ping

	s = strings.Replace(s, `"`, "", -1)

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return Ping{}, err
	}

	err = json.Unmarshal(decoded, &p)
	if err != nil {
		return Ping{}, err
	}

	return p, nil
}

func unixToTime(stamp int64) time.Time {
	return time.UnixMilli(stamp)
}

func timeToUnix(t time.Time) int64 {
	return t.Unix() * 1000
}
