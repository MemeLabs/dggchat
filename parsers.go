package dggchat

import (
	"encoding/json"
	"time"
)

func parseMessage(s string) (*Message, error) {
	var m message
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}

	message := &Message{
		Sender:    m.Nick,
		Timestamp: time.Unix(m.Timestamp/1000, 0),
		Message:   m.Data,
		Features:  m.Features,
	}

	return message, nil
}

func parseNames(s string) (*namesMessage, error) {
	var nm namesMessage
	err := json.Unmarshal([]byte(s), &nm)
	if err != nil {
		return nil, err
	}

	return &nm, nil
}

func parseRoomAction(s string) (*roomAction, error) {
	var ra roomAction

	err := json.Unmarshal([]byte(s), &ra)
	if err != nil {
		return nil, err
	}

	return &ra, nil
}
