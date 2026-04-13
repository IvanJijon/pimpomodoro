package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

func newTestModel() Model {
	return NewModel(AppConfig{
		Session: session.DefaultConfig(),
		Callbacks: Callbacks{
			PlayAlarm:  func() {},
			SendNotify: func(_, _ string) {},
		},
		ConfirmEnabled: true,
	})
}

func TestUpdateWindowSize(t *testing.T) {
	tests := []struct {
		name       string
		msg        tea.WindowSizeMsg
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "stores terminal dimensions on WindowSizeMsg",
			msg:        tea.WindowSizeMsg{Width: 120, Height: 40},
			wantWidth:  120,
			wantHeight: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()

			updated, _ := m.Update(tt.msg)
			model, ok := updated.(Model)
			if !ok {
				t.Fatal("Update did not return a Model")
			}

			if model.width != tt.wantWidth {
				t.Errorf("width = %d, want %d", model.width, tt.wantWidth)
			}
			if model.height != tt.wantHeight {
				t.Errorf("height = %d, want %d", model.height, tt.wantHeight)
			}
		})
	}
}
