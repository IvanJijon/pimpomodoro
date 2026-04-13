package theme

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) string
		wantTheme Theme
	}{
		{
			name: "loads theme from valid YAML file",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "theme.yaml")
				content := "work: \"#FF0000\"\nshort-break: \"#00FF00\"\nlong-break: \"#0000FF\"\npaused: \"#FFFFFF\"\nsubtle: \"#000000\"\n"
				err := os.WriteFile(path, []byte(content), 0644)
				assert.NoError(t, err)
				return path
			},
			wantTheme: Theme{
				Work:       "#FF0000",
				ShortBreak: "#00FF00",
				LongBreak:  "#0000FF",
				Paused:     "#FFFFFF",
				Subtle:     "#000000",
			},
		},
		{
			name: "returns default theme when file does not exist",
			setup: func(t *testing.T) string {
				t.Helper()
				return "/nonexistent/path/theme.yaml"
			},
			wantTheme: DefaultTheme(),
		},
		{
			name: "returns default theme when file is invalid YAML",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "theme.yaml")
				err := os.WriteFile(path, []byte(":::not valid yaml"), 0644)
				assert.NoError(t, err)
				return path
			},
			wantTheme: DefaultTheme(),
		},
		{
			name: "partial YAML fills missing fields with defaults",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				path := filepath.Join(dir, "theme.yaml")
				content := "work: \"#FF0000\"\n"
				err := os.WriteFile(path, []byte(content), 0644)
				assert.NoError(t, err)
				return path
			},
			wantTheme: Theme{
				Work:       "#FF0000",
				ShortBreak: "#40E0D0",
				LongBreak:  "#1E3A5F",
				Paused:     "#FFD700",
				Subtle:     "#666666",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup(t)
			got := LoadFromFile(path)

			if got != tt.wantTheme {
				t.Errorf("LoadFromFile() = %+v, want %+v", got, tt.wantTheme)
			}
		})
	}
}
