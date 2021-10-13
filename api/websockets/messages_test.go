package websockets

import (
	"testing"
)

const DealPayload string = "{\"identifier\":\"{\\\"channel\\\":\\\"DealsChannel\\\",\\\"users\\\":[{\\\"api_key\\\":\\\"34ee791264c34a2e9a16050674328065c48e8cea02384de68ee5ca9452f4cd73\\\",\\\"signature\\\":\\\"e0ece8a0090097bfbac6e2bc9625638dfd7dc9e99850cc92f4209ebc11de897d\\\"}]}\",\"message\":{\"id\":889690387,\"type\":\"Deal\",\"bot_id\":6093127,\"completed_safety_orders_count\":0,\"completed_manual_safety_orders_count\":0,\"pair\":\"USDT_BTC\",\"status\":\"base_order_placed\"}}"

func TestDealsMessage_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "Clean Path",
			payload: DealPayload,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DealsMessage{}
			if err := d.UnmarshalJSON([]byte(tt.payload)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessage_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "Clean Path - Deal",
			payload: DealPayload,
			wantErr: false,
		},
		{
			name:    "Clean Path - Confirm Subscription",
			payload: "{\"identifier\":\"{\\\"channel\\\":\\\"DealsChannel\\\",\\\"users\\\":[{\\\"api_key\\\":\\\"a1b2c3d4e5\\\",\\\"signature\\\":\\\"a1b2c3d4e5\\\"}]}\",\"type\":\"confirm_subscription\"}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Message{}
			if err := d.UnmarshalJSON([]byte(tt.payload)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPingMessage_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "Clean Path - Ping",
			payload: "{\"type\":\"ping\",\"message\":1633969206}",
			wantErr: false,
		},
		{
			name:    "Non Ping Message",
			payload: DealPayload,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &PingMessage{}
			if err := d.UnmarshalJSON([]byte(tt.payload)); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
