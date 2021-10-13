package websockets

import (
	"encoding/json"
	"errors"
	"github.com/jslowik/commacloner/api/dobjs"
)

type Message struct {
	Type       string                 `json:"type,omitempty"`
	Identifier string                 `json:"identifier,omitempty"`
	Command    string                 `json:"command,omitempty"`
	Message    map[string]interface{} `json:"message,omitempty"`
}

func (d *Message) UnmarshalJSON(data []byte) error {

	m := make(map[string]interface{})
	if e := json.Unmarshal(data, &m); e != nil {
		return e
	}

	if m["type"] != nil {
		d.Type = m["type"].(string)
	}
	if m["identifier"] != nil {
		d.Identifier = m["identifier"].(string)
	}

	details := m["message"]
	// convert map to json
	jsonString, _ := json.Marshal(details)

	// Determine if a deal
	s := dobjs.Deal{}
	if e := json.Unmarshal(jsonString, &s); e == nil && s.Type != "" {
		d.Type = s.Type
	}

	if e := json.Unmarshal(jsonString, &d.Message); e != nil {
		return e
	}
	return nil
}

type DealsMessage struct {
	Message
	Details dobjs.Deal `json:"message"`
}

type PingMessage struct {
	Message
	Time float64 `json:"message,omitempty"`
}

func (d *PingMessage) UnmarshalJSON(data []byte) error {

	m := make(map[string]interface{})
	if e := json.Unmarshal(data, &m); e != nil {
		return e
	}

	if m["type"] != nil {
		d.Type = m["type"].(string)
	}
	if m["identifier"] != nil {
		d.Identifier = m["identifier"].(string)
	}

	if d.Type != "ping" {
		return errors.New("not a ping message")
	}

	d.Time = m["message"].(float64)
	return nil
}

func (d *DealsMessage) UnmarshalJSON(data []byte) error {

	m := make(map[string]interface{})
	if e := json.Unmarshal(data, &m); e != nil {
		return e
	}

	if m["identifier"] != nil {
		d.Identifier = m["identifier"].(string)
	}

	details := m["message"]
	// convert details to json
	jsonString, _ := json.Marshal(details)

	// convert json to struct
	s := dobjs.Deal{}
	if e := json.Unmarshal(jsonString, &s); e != nil {
		return e
	}
	d.Details = s
	d.Type = s.Type
	return nil
}
