package tui

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/IvanJijon/pimpomodoro/session"
	"github.com/IvanJijon/pimpomodoro/task"
)

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
			name: "pressing d marks in-progress task at cursor as done",
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
			name: "pressing d marks pending task at cursor as done",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("First", 1))
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantStatus: task.Done,
		},
		{
			name: "pressing d on done task unmarks it to pending",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				t1 := task.NewTask("First", 1)
				m.taskList.Add(t1)
				m.taskList.MarkTaskDone(t1)
				m.taskCursor = 0
			},
			msg:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			wantStatus: task.Pending,
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

func TestTaskEdit(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*Model)
		msg          tea.Msg
		wantViewMode ViewMode
		wantName     string
		wantEstimate int
	}{
		{
			name: "pressing e enters edit mode with current values",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
				m.taskList.Add(task.NewTask("Write tests", 3))
				m.taskCursor = 0
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},
			wantViewMode: ModeTaskEdit,
			wantName:     "Write tests",
			wantEstimate: 3,
		},
		{
			name: "pressing e on empty list is a no-op",
			setup: func(m *Model) {
				m.viewMode = ModeTaskList
			},
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},
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

			if model.viewMode != tt.wantViewMode {
				t.Errorf("viewMode = %v, want %v", model.viewMode, tt.wantViewMode)
			}
			if tt.wantViewMode == ModeTaskEdit {
				if model.taskNameInput.Value() != tt.wantName {
					t.Errorf("taskNameInput = %q, want %q", model.taskNameInput.Value(), tt.wantName)
				}
				if model.taskEstimateInput.Value() != fmt.Sprintf("%d", tt.wantEstimate) {
					t.Errorf("taskEstimateInput = %q, want %q", model.taskEstimateInput.Value(), fmt.Sprintf("%d", tt.wantEstimate))
				}
			}
		})
	}
}

func TestTaskEditConfirm(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*Model)
		msg          tea.Msg
		wantViewMode ViewMode
		wantName     string
		wantEstimate int
	}{
		{
			name: "pressing enter confirms edit",
			setup: func(m *Model) {
				m.viewMode = ModeTaskEdit
				m.taskList.Add(task.NewTask("Old name", 3))
				m.taskCursor = 0
				m.taskNameInput.SetValue("New name")
				m.taskEstimateInput.SetValue("5")
			},
			msg:          tea.KeyMsg{Type: tea.KeyEnter},
			wantViewMode: ModeTaskList,
			wantName:     "New name",
			wantEstimate: 5,
		},
		{
			name: "pressing esc cancels edit",
			setup: func(m *Model) {
				m.viewMode = ModeTaskEdit
				m.taskList.Add(task.NewTask("Old name", 3))
				m.taskCursor = 0
				m.taskNameInput.SetValue("New name")
				m.taskEstimateInput.SetValue("5")
			},
			msg:          tea.KeyMsg{Type: tea.KeyEsc},
			wantViewMode: ModeTaskList,
			wantName:     "Old name",
			wantEstimate: 3,
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
			tsk := model.taskList.Tasks()[0]
			if tsk.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tsk.Name, tt.wantName)
			}
			if tsk.EstimatedPomos != tt.wantEstimate {
				t.Errorf("EstimatedPomos = %d, want %d", tsk.EstimatedPomos, tt.wantEstimate)
			}
		})
	}
}
