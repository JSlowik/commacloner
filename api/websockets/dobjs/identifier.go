package dobjs

// Identifier the websocket identifier
type Identifier struct {
	Channel string `json:"channel"`
	Users   []User `json:"users"`
}

// User the user/api subscribing to the websocket
type User struct {
	APIKey    string `json:"api_key"`
	Signature string `json:"signature"`
}

// IdentifierMessage the message to be sent to the websocket url to subscribe to a stream
type IdentifierMessage struct {
	Identifier string `json:"identifier,omitempty"`
	Command    string `json:"command,omitempty"`
	Type       string `json:"type,omitempty"`
}
