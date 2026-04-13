package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/task"
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

func TestVisualAlert(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*Model)
		msg            tea.Msg
		wantAlerting   bool
		wantBlinkState bool
	}{
		{
			name: "timer expiry sets alerting when visual alert enabled",
			setup: func(m *Model) {
				m.visualAlert = true
				m.session.CurrentPhase = session.Work
				m.remainingTime = 0
				m.running = true
			},
			msg:          TickMsg{},
			wantAlerting: true,
		},
		{
			name: "timer expiry does not set alerting when visual alert disabled",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 0
				m.running = true
			},
			msg:          TickMsg{},
			wantAlerting: false,
		},
		{
			name: "any keypress clears alerting",
			setup: func(m *Model) {
				m.alerting = true
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}},
			wantAlerting: false,
		},
		{
			name: "blink msg toggles blinkState while alerting",
			setup: func(m *Model) {
				m.alerting = true
				m.blinkState = false
			},
			msg:            BlinkMsg{},
			wantAlerting:   true,
			wantBlinkState: true,
		},
		{
			name: "blink msg is ignored when not alerting",
			setup: func(m *Model) {
				m.alerting = false
				m.blinkState = false
			},
			msg:            BlinkMsg{},
			wantAlerting:   false,
			wantBlinkState: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()
			if tt.setup != nil {
				tt.setup(&m)
			}

			updated, _ := m.Update(tt.msg)
			model := updated.(Model)

			if model.alerting != tt.wantAlerting {
				t.Errorf("alerting = %v, want %v", model.alerting, tt.wantAlerting)
			}
			if model.blinkState != tt.wantBlinkState {
				t.Errorf("blinkState = %v, want %v", model.blinkState, tt.wantBlinkState)
			}
		})
	}
}

func TestWorkPhaseCompletionIncrementsWIPPomos(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*Model)
		msg        tea.Msg
		wantActual int
	}{
		{
			name: "Work phase completion increments WIP task actual pomodoros",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 0
				m.running = true
				wip := task.NewTask("Write tests", 3)
				m.taskList.Add(wip)
				m.taskList.SelectWIP(wip)
			},
			msg:        TickMsg{},
			wantActual: 1,
		},
		{
			name: "ShortBreak phase completion does not increment WIP task",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.ShortBreak
				m.remainingTime = 0
				m.running = true
				wip := task.NewTask("Write tests", 3)
				m.taskList.Add(wip)
				m.taskList.SelectWIP(wip)
			},
			msg:        TickMsg{},
			wantActual: 0,
		},
		{
			name: "Work phase completion with no WIP task does not panic",
			setup: func(m *Model) {
				m.session.CurrentPhase = session.Work
				m.remainingTime = 0
				m.running = true
			},
			msg:        TickMsg{},
			wantActual: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()
			if tt.setup != nil {
				tt.setup(&m)
			}

			updated, _ := m.Update(tt.msg)
			model := updated.(Model)

			wip := model.taskList.CurrentWIP()
			if tt.wantActual == -1 {
				if wip != nil {
					t.Errorf("CurrentWIP() = %q, want nil", wip.Name)
				}
				return
			}
			if wip == nil {
				t.Fatal("CurrentWIP() = nil, want a task")
			}
			if wip.ActualPomos != tt.wantActual {
				t.Errorf("ActualPomos = %d, want %d", wip.ActualPomos, tt.wantActual)
			}
		})
	}
}
