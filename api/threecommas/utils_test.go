package threecommas

import "testing"

func TestComputeSignature(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		secret string
		want   string
	}{
		{
			name:   "clean path",
			path:   "/deals",
			secret: "s0m3s3cr3t!!",
			want:   "0a77586521ce9d268f87e6d3bcf5a3c0995481c37dce4502914d07f61562f57f",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeSignature(tt.path, tt.secret); got != tt.want {
				t.Errorf("ComputeSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
