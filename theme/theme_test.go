package theme

import "testing"

func TestDefaultTheme(t *testing.T) {
	th := DefaultTheme()

	tests := []struct {
		name      string
		got       string
		wantColor string
	}{
		{name: "Work", got: th.Work, wantColor: "#722F37"},
		{name: "ShortBreak", got: th.ShortBreak, wantColor: "#40E0D0"},
		{name: "LongBreak", got: th.LongBreak, wantColor: "#1E3A5F"},
		{name: "Paused", got: th.Paused, wantColor: "#FFD700"},
		{name: "Subtle", got: th.Subtle, wantColor: "#666666"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.wantColor {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.wantColor)
			}
		})
	}
}
