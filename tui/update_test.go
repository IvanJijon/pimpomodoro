package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

func newTestModel() Model {
	cb := Callbacks{
		PlayAlarm:  func() {},            // no-op
		SendNotify: func(_, _ string) {}, // no-op
	}
	return NewModel(session.DefaultConfig(), cb)
}

func TestUpdateTick(t *testing.T) {
	tests := []struct {
		name              string
		setup             func(*Model)
		msg               tea.Msg
		wantPhase         session.Phase
		wantRemainingTime time.Duration
		wantRunning       bool
	}{
		{
			name: "tick decrements remaining time by one second",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 25 * time.Minute
				m.running = true
			},
			msg:               TickMsg{},
			wantPhase:         session.Work,
			wantRemainingTime: 24*time.Minute + 59*time.Second,
			wantRunning:       true,
		},
		{
			name: "tick at zero stops timer and transitions to next phase",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 0
				m.running = true
			},
			msg:               TickMsg{},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 5 * time.Minute,
			wantRunning:       false,
		},
		{
			name: "tick at zero on last pomodoro transitions to LongBreak",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = m.session.Rounds
				m.remainingTime = 0
				m.running = true
			},
			msg:               TickMsg{},
			wantPhase:         session.LongBreak,
			wantRemainingTime: 15 * time.Minute,
			wantRunning:       false,
		},
		{
			name: "tick when not running is a no-op",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.remainingTime = 5 * time.Minute
				m.running = false
			},
			msg:               TickMsg{},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 5 * time.Minute,
			wantRunning:       false,
		},
		{
			name: "stale tick from old loop is ignored",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 25 * time.Minute
				m.running = true
				m.tickID = 2
			},
			msg:               TickMsg{id: 1},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()
			if tt.setup != nil {
				tt.setup(&m)
			}

			updated, _ := m.Update(tt.msg)
			model, ok := updated.(Model)
			if !ok {
				t.Fatal("Update did not return a Model")
			}

			if model.session.CurrentPhase != tt.wantPhase {
				t.Errorf("CurrentPhase = %v, want %v", model.session.CurrentPhase, tt.wantPhase)
			}
			if model.remainingTime != tt.wantRemainingTime {
				t.Errorf("remainingTime = %v, want %v", model.remainingTime, tt.wantRemainingTime)
			}
			if model.running != tt.wantRunning {
				t.Errorf("running = %v, want %v", model.running, tt.wantRunning)
			}
		})
	}
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

func TestUpdateKeyMsg(t *testing.T) {
	tests := []struct {
		name              string
		setup             func(*Model)
		msg               tea.Msg
		wantPhase         session.Phase
		wantRemainingTime time.Duration
		wantRunning       bool
		wantViewMode      ViewMode
	}{
		{
			name:              "pressing s while Idle starts Work phase",
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       true,
		},
		{
			name: "pressing s while paused resumes without resetting time",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.running = false
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
		},
		{
			name: "pressing s while running pauses without resetting time",
			setup: func(m *Model) {
				updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
				*m = updated.(Model)
				updated, _ = m.Update(TickMsg{id: m.tickID})
				*m = updated.(Model)
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
			wantPhase:         session.Work,
			wantRemainingTime: 24*time.Minute + 59*time.Second,
			wantRunning:       false,
		},
		{
			name: "pressing r shows reset confirmation",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.running = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeResetConfirm,
		},
		{
			name: "pressing y during reset confirmation resets timer",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeResetConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing n during reset confirmation cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeResetConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing n shows skip confirmation",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = 1
				m.remainingTime = 12 * time.Minute
				m.running = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeSkipConfirm,
		},
		{
			name: "pressing y during confirm skip skips to next phase",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = 1
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeSkipConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 5 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing n during confirm skip cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = 1
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeSkipConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name:              "pressing n while Idle is a no-op",
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:         session.Idle,
			wantRemainingTime: 0,
			wantRunning:       false,
		},
		{
			name: "pressing b goes to previous phase",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.session.CurrentPomodoro = 2
				m.remainingTime = 3 * time.Minute
				m.running = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       false,
		},
		{
			name:              "pressing b while Idle is a no-op",
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}},
			wantPhase:         session.Idle,
			wantRemainingTime: 0,
			wantRunning:       false,
		},
		{
			name:         "pressing ? enters help mode",
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
			wantViewMode: ModeHelp,
		},
		{
			name: "pressing ? again exits help mode",
			setup: func(m *Model) {
				m.viewMode = ModeHelp
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
			wantViewMode: ModeNormal,
		},
		{
			name: "pressing q shows quit confirmation",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.running = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeQuitConfirm,
		},
		{
			name: "pressing n during quit confirmation cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeQuitConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing n during quit confirmation while Idle stays Idle",
			setup: func(m *Model) {
				m.viewMode = ModeQuitConfirm
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:    session.Idle,
			wantRunning:  false,
			wantViewMode: ModeNormal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()
			if tt.setup != nil {
				tt.setup(&m)
			}

			updated, _ := m.Update(tt.msg)
			model, ok := updated.(Model)
			if !ok {
				t.Fatal("Update did not return a Model")
			}

			if model.session.CurrentPhase != tt.wantPhase {
				t.Errorf("CurrentPhase = %v, want %v", model.session.CurrentPhase, tt.wantPhase)
			}
			if model.remainingTime != tt.wantRemainingTime {
				t.Errorf("remainingTime = %v, want %v", model.remainingTime, tt.wantRemainingTime)
			}
			if model.running != tt.wantRunning {
				t.Errorf("running = %v, want %v", model.running, tt.wantRunning)
			}
			if model.viewMode != tt.wantViewMode {
				t.Errorf("viewMode = %v, want %v", model.viewMode, tt.wantViewMode)
			}
		})
	}
}
