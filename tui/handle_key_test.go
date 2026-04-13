package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
)

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
			name: "pressing r with confirmations disabled resets directly",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.running = true
				m.confirmEnabled = false
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeNormal,
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
			name: "pressing n with confirmations disabled skips directly",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = 1
				m.remainingTime = 12 * time.Minute
				m.running = true
				m.confirmEnabled = false
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 5 * time.Minute,
			wantRunning:       false,
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
			name: "pressing b shows previous phase confirmation",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.session.CurrentPomodoro = 2
				m.remainingTime = 3 * time.Minute
				m.running = true
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 3 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModePreviousConfirm,
		},
		{
			name: "pressing y during previous confirmation goes to previous phase",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.session.CurrentPomodoro = 2
				m.remainingTime = 3 * time.Minute
				m.viewMode = ModePreviousConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing n during previous confirmation cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.session.CurrentPomodoro = 2
				m.remainingTime = 3 * time.Minute
				m.viewMode = ModePreviousConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 3 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing b with confirmations disabled goes directly to previous phase",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.session.CurrentPomodoro = 2
				m.remainingTime = 3 * time.Minute
				m.running = true
				m.confirmEnabled = false
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}},
			wantPhase:         session.Work,
			wantRemainingTime: 25 * time.Minute,
			wantRunning:       false,
			wantViewMode:      ModeNormal,
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
			name:         "pressing t in normal mode enters task list mode",
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}},
			wantViewMode: ModeTaskList,
		},
		{
			name: "pressing t in task list mode returns to normal mode",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}},
			wantViewMode: ModeNormal,
		},
		{
			name: "pressing a in task list mode enters task add mode",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantViewMode: ModeTaskAdd,
		},
		{
			name: "pressing esc in task add mode returns to task list",
			setup: func(m *Model) {
				m.viewMode = ModeTaskAdd
			},
			msg:          tea.KeyMsg{Type: tea.KeyEsc},
			wantViewMode: ModeTaskList,
		},
		{
			name: "pressing esc in help mode returns to normal",
			setup: func(m *Model) {
				m.viewMode = ModeHelp
			},
			msg:          tea.KeyMsg{Type: tea.KeyEsc},
			wantViewMode: ModeNormal,
		},
		{
			name: "pressing esc in task list mode returns to normal",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
			},
			msg:          tea.KeyMsg{Type: tea.KeyEsc},
			wantViewMode: ModeNormal,
		},
		{
			name: "pressing esc in reset confirmation cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeResetConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyEsc},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing esc in skip confirmation cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.session.CurrentPomodoro = 1
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeSkipConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyEsc},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing esc in previous confirmation cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.session.CurrentPomodoro = 2
				m.remainingTime = 3 * time.Minute
				m.viewMode = ModePreviousConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyEsc},
			wantPhase:         session.ShortBreak,
			wantRemainingTime: 3 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing esc in quit confirmation cancels and resumes",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.viewMode = ModeQuitConfirm
			},
			msg:               tea.KeyMsg{Type: tea.KeyEsc},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
		},
		{
			name: "pressing q from task list shows quit confirmation",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantViewMode: ModeQuitConfirm,
		},
		{
			name: "pressing q from help shows quit confirmation",
			setup: func(m *Model) {
				m.viewMode = ModeHelp
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantViewMode: ModeQuitConfirm,
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
			name: "pressing q with confirmations disabled quits directly",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				m.running = true
				m.confirmEnabled = false
			},
			msg:               tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantPhase:         session.Work,
			wantRemainingTime: 12 * time.Minute,
			wantRunning:       true,
			wantViewMode:      ModeNormal,
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
