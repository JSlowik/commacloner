package api

import (
	"encoding/json"
	"errors"
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
	s := DealDetails{}
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
	Details DealDetails `json:"message"`
}

type DealDetails struct {
	ID                               int    `json:"id"`
	Type                             string `json:"type"`
	BotID                            int    `json:"bot_id"`
	CompletedSafetyOrdersCount       int    `json:"completed_safety_orders_count"`
	CompletedManualSafetyOrdersCount int    `json:"completed_manual_safety_orders_count"`
	Pair                             string `json:"pair"`
	Status                           string `json:"status"`
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
	s := DealDetails{}
	if e := json.Unmarshal(jsonString, &s); e != nil {
		return e
	}
	d.Details = s
	d.Type = s.Type
	return nil
}
