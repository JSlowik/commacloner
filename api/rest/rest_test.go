package rest

import "testing"

func Test_generateQuery(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		queryParameters map[string]string
		want            string
	}{
		{
			name: "generate new deal test",
			path: "https://api.3commas.io/public/api/ver1/bots/1234/start_new_deal",
			queryParameters: map[string]string{
				"pair":               "BTC_USD",
				"skip_signal_checks": "true",
			},
			want: "https://api.3commas.io/public/api/ver1/bots/1234/start_new_deal?pair=BTC_USD&skip_signal_checks=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateQuery(tt.path, tt.queryParameters); got.String() != tt.want {
				t.Errorf("generateQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
