package config

import (
	"testing"
)

func TestCgearHomeFromExecutable(t *testing.T) {
	tests := []struct {
		name       string
		executable string
		want       string
	}{
		{
			name:       "installed layout",
			executable: `C:\Tools\cgear\bin\cgear.exe`,
			want:       `C:\Tools\cgear`,
		},
		{
			name:       "standalone executable",
			executable: `C:\Tools\cgear\cgear.exe`,
			want:       `C:\Tools\cgear`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cgearHomeFromExecutable(tt.executable); got != tt.want {
				t.Fatalf("cgearHomeFromExecutable() = %q, want %q", got, tt.want)
			}
		})
	}
}
