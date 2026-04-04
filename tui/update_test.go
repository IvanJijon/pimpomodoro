package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

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
			m := NewModel()
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

func TestUpdateKeyMsg(t *testing.T) {
	tests := []struct {
		name              string
		setup             func(*Model)
		msg               tea.Msg
		wantPhase         session.Phase
		wantRemainingTime time.Duration
		wantRunning       bool
		wantShowHelp      bool
		wantConfirmSkip   bool
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
			name: "pressing s while running does not reset time",
			setup: func(m *Model) {
				// First press: start the session
				updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
				*m = updated.(Model)
				// Simulate some time passing
				updated, _ = m.Update(TickMsg{id: m.tickID})
				*m = updated.(Model)
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
			wantPhase:         session.Work,
			wantRemainingTime: 24*time.Minute + 59*time.Second,
			wantRunning:       true,
		},
		{
			name: "pressing p pauses the timer",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 20 * time.Minute
				m.running = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}},
			wantPhase:         session.Work,
			wantRemainingTime: 20 * time.Minute,
			wantRunning:       false,
		},
		{
			name: "pressing r resets current phase timer",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.running = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       false,
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
			wantConfirmSkip:   true,
		},
		{
			name: "pressing y during confirm skip skips to next phase",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = 1
				m.remainingTime = 12 * time.Minute
				m.showSkipConfirm = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 5 * time.Minute,
			wantRunning:       false,
			wantConfirmSkip:   false,
		},
		{
			name: "pressing x during confirm skip cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = 1
				m.remainingTime = 12 * time.Minute
				m.showSkipConfirm = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantConfirmSkip:   false,
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
			name:         "pressing ? toggles help on",
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
			wantShowHelp: true,
		},
		{
			name: "pressing ? again toggles help off",
			setup: func(m *Model) {
				m.showHelp = true
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
			wantShowHelp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
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
			if model.showHelp != tt.wantShowHelp {
				t.Errorf("showHelp = %v, want %v", model.showHelp, tt.wantShowHelp)
			}
			if model.showSkipConfirm != tt.wantConfirmSkip {
				t.Errorf("showSkipConfirm = %v, want %v", model.showSkipConfirm, tt.wantConfirmSkip)
			}
		})
	}
}
