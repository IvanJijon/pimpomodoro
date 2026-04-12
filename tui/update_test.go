package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/task"
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

func TestTaskListNavigation(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*Model)
		msg            tea.Msg
		wantTaskCursor int
	}{
		{
			name: "pressing down moves cursor down",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
			},
			msg:            tea.KeyMsg{Type: tea.KeyDown},
			wantTaskCursor: 1,
		},
		{
			name: "pressing up from second item moves cursor up",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
				m.taskCursor = 1
			},
			msg:            tea.KeyMsg{Type: tea.KeyUp},
			wantTaskCursor: 0,
		},
		{
			name: "pressing down moves cursor down",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
			},
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			wantTaskCursor: 1,
		},
		{
			name: "pressing up from second item moves cursor up",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
				m.taskCursor = 1
			},
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			wantTaskCursor: 0,
		},
		{
			name: "cursor does not go below last item",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
				m.taskCursor = 1
			},
			msg:            tea.KeyMsg{Type: tea.KeyDown},
			wantTaskCursor: 1,
		},
		{
			name: "cursor does not go above first item",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
				m.taskCursor = 0
			},
			msg:            tea.KeyMsg{Type: tea.KeyUp},
			wantTaskCursor: 0,
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

			if model.taskCursor != tt.wantTaskCursor {
				t.Errorf("taskCursor = %d, want %d", model.taskCursor, tt.wantTaskCursor)
			}
		})
	}
}

func TestTaskListSelectWIPFromCursor(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*Model)
		msg         tea.Msg
		wantWIPName string
	}{
		{
			name: "pressing enter selects task at cursor as WIP",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
				m.taskCursor = 1
			},
			msg:         tea.KeyMsg{Type: tea.KeyEnter},
			wantWIPName: "Second",
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
			if wip == nil {
				t.Fatal("CurrentWIP() = nil, want a task")
			}
			if wip.Name != tt.wantWIPName {
				t.Errorf("CurrentWIP().Name = %q, want %q", wip.Name, tt.wantWIPName)
			}
		})
	}
}

func TestTaskListEnterOnEmptyList(t *testing.T) {
	m := newTestModel()
	m.viewMode = ModeTaskList

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model := updated.(Model)

	if wip := model.taskList.CurrentWIP(); wip != nil {
		t.Errorf("CurrentWIP() = %q, want nil", wip.Name)
	}
}

func TestTaskListMarkDoneFromCursor(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*Model)
		msg        tea.Msg
		wantStatus task.Status
	}{
		{
			name: "pressing d marks task at cursor as done",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				t1 := task.NewTask("First", 1)
				m.taskList.Add(t1)
				m.taskList.SelectWIP(t1)
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantStatus: task.Done,
		},
		{
			name: "pressing d on empty list is a no-op",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantStatus: 0,
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

			tasks := model.taskList.Tasks()
			if len(tasks) == 0 {
				return
			}
			if tasks[0].Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", tasks[0].Status, tt.wantStatus)
			}
		})
	}
}

func TestTaskListRemoveFromCursor(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*Model)
		msg        tea.Msg
		wantCount  int
		wantCursor int
	}{
		{
			name: "pressing x removes task at cursor",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
				m.taskCursor = 0
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
			wantCount:  1,
			wantCursor: 0,
		},
		{
			name: "pressing x on empty list is a no-op",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
			wantCount:  0,
			wantCursor: 0,
		},
		{
			name: "pressing x on last item adjusts cursor to previous",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
				m.taskList.Add(task.NewTask("Second", 2))
				m.taskCursor = 1
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
			wantCount:  1,
			wantCursor: 0,
		},
		{
			name: "pressing x on only remaining task resets cursor to zero",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("Only", 1))
				m.taskCursor = 0
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
			wantCount:  0,
			wantCursor: 0,
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

			if model.taskList.Len() != tt.wantCount {
				t.Errorf("Len() = %d, want %d", model.taskList.Len(), tt.wantCount)
			}
			if model.taskCursor != tt.wantCursor {
				t.Errorf("taskCursor = %d, want %d", model.taskCursor, tt.wantCursor)
			}
		})
	}
}

func TestTaskAddConfirm(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*Model)
		msg          tea.Msg
		wantCount    int
		wantTaskName string
		wantEstimate int
		wantViewMode ViewMode
	}{
		{
			name: "pressing enter in task add mode adds task and returns to task list",
			setup: func(m *Model) {
				m.viewMode = ModeTaskAdd
				m.taskNameInput.SetValue("Write tests")
				m.taskEstimateInput.SetValue("3")
			},
			msg:          tea.KeyMsg{Type: tea.KeyEnter},
			wantCount:    1,
			wantTaskName: "Write tests",
			wantEstimate: 3,
			wantViewMode: ModeTaskList,
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

			if model.taskList.Len() != tt.wantCount {
				t.Errorf("Len() = %d, want %d", model.taskList.Len(), tt.wantCount)
			}
			if model.viewMode != tt.wantViewMode {
				t.Errorf("viewMode = %v, want %v", model.viewMode, tt.wantViewMode)
			}
			if model.taskList.Len() > 0 {
				tsk := model.taskList.Tasks()[0]
				if tsk.Name != tt.wantTaskName {
					t.Errorf("Name = %q, want %q", tsk.Name, tt.wantTaskName)
				}
				if tsk.EstimatedPomos != tt.wantEstimate {
					t.Errorf("EstimatedPomos = %d, want %d", tsk.EstimatedPomos, tt.wantEstimate)
				}
			}
		})
	}
}

func TestTaskAddInvalidEstimate(t *testing.T) {
	tests := []struct {
		name         string
		estimate     string
		wantEstimate int
	}{
		{
			name:         "empty estimate defaults to 1",
			estimate:     "",
			wantEstimate: 1,
		},
		{
			name:         "non-numeric estimate defaults to 1",
			estimate:     "abc",
			wantEstimate: 1,
		},
		{
			name:         "negative estimate defaults to 1",
			estimate:     "-3",
			wantEstimate: 1,
		},
		{
			name:         "zero estimate defaults to 1",
			estimate:     "0",
			wantEstimate: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestModel()
			m.viewMode = ModeTaskAdd
			m.taskNameInput.SetValue("Write tests")
			m.taskEstimateInput.SetValue(tt.estimate)

			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model := updated.(Model)

			if model.taskList.Len() != 1 {
				t.Fatalf("Len() = %d, want 1", model.taskList.Len())
			}
			if model.taskList.Tasks()[0].EstimatedPomos != tt.wantEstimate {
				t.Errorf("EstimatedPomos = %d, want %d", model.taskList.Tasks()[0].EstimatedPomos, tt.wantEstimate)
			}
		})
	}
}

func TestTaskAddEmptyName(t *testing.T) {
	m := newTestModel()
	m.viewMode = ModeTaskAdd
	m.taskNameInput.SetValue("")
	m.taskEstimateInput.SetValue("3")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model := updated.(Model)

	if model.taskList.Len() != 0 {
		t.Errorf("Len() = %d, want 0", model.taskList.Len())
	}
	if model.viewMode != ModeTaskAdd {
		t.Errorf("viewMode = %v, want %v", model.viewMode, ModeTaskAdd)
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

func TestTaskSwitchConfirmation(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*Model)
		msg          tea.Msg
		wantViewMode ViewMode
		wantWIPName  string
	}{
		{
			name: "pressing enter while timer running shows switch confirmation",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.running = true
				m.session.CurrentPhase = session.Work
				first := task.NewTask("First", 1)
				second := task.NewTask("Second", 2)
				m.taskList.Add(first)
				m.taskList.Add(second)
				m.taskList.SelectWIP(first)
				m.taskCursor = 1
			},
			msg:          tea.KeyMsg{Type: tea.KeyEnter},
			wantViewMode: ModeSwitchTaskConfirm,
			wantWIPName:  "First",
		},
		{
			name: "pressing enter while timer not running switches directly when Idle",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.running = false
				first := task.NewTask("First", 1)
				second := task.NewTask("Second", 2)
				m.taskList.Add(first)
				m.taskList.Add(second)
				m.taskList.SelectWIP(first)
				m.taskCursor = 1
			},
			msg:          tea.KeyMsg{Type: tea.KeyEnter},
			wantViewMode: ModeTaskList,
			wantWIPName:  "Second",
		},
		{
			name: "pressing enter while timer paused shows switch confirmation",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.running = false
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				first := task.NewTask("First", 1)
				second := task.NewTask("Second", 2)
				m.taskList.Add(first)
				m.taskList.Add(second)
				m.taskList.SelectWIP(first)
				m.taskCursor = 1
			},
			msg:          tea.KeyMsg{Type: tea.KeyEnter},
			wantViewMode: ModeSwitchTaskConfirm,
			wantWIPName:  "First",
		},
		{
			name: "pressing enter on already WIP task while running is a no-op",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.running = true
				m.session.CurrentPhase = session.Work
				first := task.NewTask("First", 1)
				m.taskList.Add(first)
				m.taskList.SelectWIP(first)
				m.taskCursor = 0
			},
			msg:          tea.KeyMsg{Type: tea.KeyEnter},
			wantViewMode: ModeTaskList,
			wantWIPName:  "First",
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

			if model.viewMode != tt.wantViewMode {
				t.Errorf("viewMode = %v, want %v", model.viewMode, tt.wantViewMode)
			}
			wip := model.taskList.CurrentWIP()
			if wip == nil {
				t.Fatal("CurrentWIP() = nil, want a task")
			}
			if wip.Name != tt.wantWIPName {
				t.Errorf("CurrentWIP().Name = %q, want %q", wip.Name, tt.wantWIPName)
			}
		})
	}
}

func TestSwitchTaskConfirmDialog(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*Model)
		msg          tea.Msg
		wantViewMode ViewMode
		wantWIPName  string
		wantRunning  bool
	}{
		{
			name: "pressing y confirms task switch and resets timer",
			setup: func(m *Model) {
				m.viewMode = ModeSwitchTaskConfirm
				m.session.CurrentPhase = session.Work
				m.remainingTime = 12 * time.Minute
				first := task.NewTask("First", 1)
				second := task.NewTask("Second", 2)
				m.taskList.Add(first)
				m.taskList.Add(second)
				m.taskList.SelectWIP(first)
				m.taskCursor = 1
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantViewMode: ModeTaskList,
			wantWIPName:  "Second",
			wantRunning:  false,
		},
		{
			name: "pressing n cancels task switch",
			setup: func(m *Model) {
				m.viewMode = ModeSwitchTaskConfirm
				m.session.CurrentPhase = session.Work
				first := task.NewTask("First", 1)
				second := task.NewTask("Second", 2)
				m.taskList.Add(first)
				m.taskList.Add(second)
				m.taskList.SelectWIP(first)
				m.taskCursor = 1
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantViewMode: ModeTaskList,
			wantWIPName:  "First",
			wantRunning:  false,
		},
		{
			name: "pressing esc cancels task switch",
			setup: func(m *Model) {
				m.viewMode = ModeSwitchTaskConfirm
				m.session.CurrentPhase = session.Work
				first := task.NewTask("First", 1)
				second := task.NewTask("Second", 2)
				m.taskList.Add(first)
				m.taskList.Add(second)
				m.taskList.SelectWIP(first)
				m.taskCursor = 1
			},
			msg:          tea.KeyMsg{Type: tea.KeyEsc},
			wantViewMode: ModeTaskList,
			wantWIPName:  "First",
			wantRunning:  false,
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

			if model.viewMode != tt.wantViewMode {
				t.Errorf("viewMode = %v, want %v", model.viewMode, tt.wantViewMode)
			}
			wip := model.taskList.CurrentWIP()
			if wip == nil {
				t.Fatal("CurrentWIP() = nil, want a task")
			}
			if wip.Name != tt.wantWIPName {
				t.Errorf("CurrentWIP().Name = %q, want %q", wip.Name, tt.wantWIPName)
			}
			if model.running != tt.wantRunning {
				t.Errorf("running = %v, want %v", model.running, tt.wantRunning)
			}
		})
	}
}
