package tui

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "formats 25 minutes as 25:00",
			duration: 25 * time.Minute,
			want:     "25:00",
		},
		{
			name:     "formats 5 minutes as 05:00",
			duration: 5 * time.Minute,
			want:     "05:00",
		},
		{
			name:     "formats 24 minutes 59 seconds as 24:59",
			duration: 24*time.Minute + 59*time.Second,
			want:     "24:59",
		},
		{
			name:     "formats zero as 00:00",
			duration: 0,
			want:     "00:00",
		},
		{
			name:     "formats 1 second as 00:01",
			duration: 1 * time.Second,
			want:     "00:01",
		},
		{
			name:     "formats 59 seconds as 00:59",
			duration: 59 * time.Second,
			want:     "00:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration() = %q, want %q", got, tt.want)
			}
		})
	}
}
