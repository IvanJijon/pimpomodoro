package tui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/theme"
)

func TestPhaseColor_theme(t *testing.T) {
	customTheme := theme.Theme{
		Work:       "#FF0000",
		ShortBreak: "#00FF00",
		LongBreak:  "#0000FF",
		Paused:     "#FFFFFF",
		Subtle:     "#000000",
	}

	tests := []struct {
		name      string
		phase     session.Phase
		wantColor lipgloss.Color
	}{
		{name: "Work uses theme Work color", phase: session.Work, wantColor: lipgloss.Color("#FF0000")},
		{name: "ShortBreak uses theme ShortBreak color", phase: session.ShortBreak, wantColor: lipgloss.Color("#00FF00")},
		{name: "LongBreak uses theme LongBreak color", phase: session.LongBreak, wantColor: lipgloss.Color("#0000FF")},
		{name: "Idle uses theme Subtle color", phase: session.Idle, wantColor: lipgloss.Color("#000000")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(AppConfig{Theme: customTheme})
			m.session.CurrentPhase = tt.phase

			got := phaseColor(m)
			if got != tt.wantColor {
				t.Errorf("phaseColor() = %q, want %q", got, tt.wantColor)
			}
		})
	}
}

func TestNewModel_theme(t *testing.T) {
	tests := []struct {
		name      string
		cfg       AppConfig
		wantTheme theme.Theme
	}{
		{
			name: "uses theme from AppConfig when provided",
			cfg: AppConfig{
				Theme: theme.Theme{
					Work:       "#FF0000",
					ShortBreak: "#00FF00",
					LongBreak:  "#0000FF",
					Paused:     "#FFFFFF",
					Subtle:     "#000000",
				},
			},
			wantTheme: theme.Theme{
				Work:       "#FF0000",
				ShortBreak: "#00FF00",
				LongBreak:  "#0000FF",
				Paused:     "#FFFFFF",
				Subtle:     "#000000",
			},
		},
		{
			name:      "defaults to DefaultTheme when zero value",
			cfg:       AppConfig{}, // zero value with no theme set
			wantTheme: theme.DefaultTheme(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(tt.cfg)

			if m.theme != tt.wantTheme {
				t.Errorf("theme = %+v, want %+v", m.theme, tt.wantTheme)
			}
		})
	}
}
