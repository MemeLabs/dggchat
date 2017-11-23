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
		Nick:     m.Nick,
		Features: m.Features,
	}

	message := Message{
		Sender:    user,
		Timestamp: unixToTime(m.Timestamp),
		Message:   m.Data,
	}

	return message, nil
}

func parseNames(s string) (namesMessage, error) {
	var nm namesMessage
	err := json.Unmarshal([]byte(s), &nm)
	if err != nil {
		return namesMessage{}, err
	}

	return nm, nil
}

func parseRoomAction(s string) (RoomAction, error) {
	var ra roomAction

	err := json.Unmarshal([]byte(s), &ra)
	if err != nil {
		return RoomAction{}, err
	}

	roomAction := RoomAction{
		User: User{
			Nick:     ra.Nick,
			Features: ra.Features,
		},
		Timestamp: unixToTime(ra.Timestamp),
	}

	return roomAction, nil
}

func parseErrorMessage(s string) string {
	return strings.Replace(s, `"`, "", -1)
}

func parsePrivateMessage(s string) (PrivateMessage, error) {
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

	return privateMessage, nil
}

func parseBroadcast(s string) (Broadcast, error) {
	var b broadcast

	err := json.Unmarshal([]byte(s), &b)
	if err != nil {
		return Broadcast{}, err
	}

	broadcast := Broadcast{
		Message:   b.Data,
		Timestamp: unixToTime(b.Timestamp),
	}

	return broadcast, nil
}

func parsePing(s string) (ping, error) {
	var p ping

	s = strings.Replace(s, `"`, "", -1)

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ping{}, err
	}

	err = json.Unmarshal(decoded, &p)
	if err != nil {
		return ping{}, err
	}

	return p, nil
}

func unixToTime(stamp int64) time.Time {
	return time.Unix(stamp/1000, 0)
}

func timeToUnix(t time.Time) int64 {
	return t.Unix() * 1000
}
