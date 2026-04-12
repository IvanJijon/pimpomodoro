package task

import "testing"

func TestNewTask(t *testing.T) {
	tests := []struct {
		name           string
		taskName       string
		estimatedPomos int
		wantStatus     Status
		wantActual     int
		wantEstimate   int
	}{
		{
			name:           "new task has pending status and zero actual pomodoros",
			taskName:       "Write tests",
			estimatedPomos: 3,
			wantStatus:     Pending,
			wantActual:     0,
			wantEstimate:   3,
		},
		{
			name:           "negative estimate defaults to 1",
			taskName:       "Write tests",
			estimatedPomos: -3,
			wantStatus:     Pending,
			wantActual:     0,
			wantEstimate:   1,
		},
		{
			name:           "zero estimate defaults to 1",
			taskName:       "Write tests",
			estimatedPomos: 0,
			wantStatus:     Pending,
			wantActual:     0,
			wantEstimate:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTask(tt.taskName, tt.estimatedPomos)

			if task.Name != tt.taskName {
				t.Errorf("Name = %q, want %q", task.Name, tt.taskName)
			}
			if task.EstimatedPomos != tt.wantEstimate {
				t.Errorf("EstimatedPomos = %d, want %d", task.EstimatedPomos, tt.wantEstimate)
			}
			if task.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
			}
			if task.ActualPomos != tt.wantActual {
				t.Errorf("ActualPomos = %d, want %d", task.ActualPomos, tt.wantActual)
			}
		})
	}
}

func TestTaskStartWork(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *Task
		wantStatus Status
	}{
		{
			name: "pending task can be started",
			setup: func() *Task {
				return NewTask("Write tests", 3)
			},
			wantStatus: InProgress,
		},
		{
			name: "starting an in-progress task is a no-op",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				task.StartWork()
				return task
			},
			wantStatus: InProgress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setup()
			task.StartWork()

			if task.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
			}
		})
	}
}

func TestTaskMarkDone(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *Task
		wantStatus Status
	}{
		{
			name: "in-progress task can be marked done",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				task.StartWork()
				return task
			},
			wantStatus: Done,
		},
		{
			name: "pending task cannot be marked done",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				return task
			},
			wantStatus: Pending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setup()
			task.MarkDone()

			if task.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
			}
		})
	}
}

func TestTaskIncreaseActualPomos(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *Task
		wantActual int
	}{
		{
			name: "increases actual pomodoros by one",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				task.StartWork()
				return task
			},
			wantActual: 1,
		},
		{
			name: "increases actual pomodoros cumulatively",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				task.StartWork()
				task.IncreaseActualPomos()
				return task
			},
			wantActual: 2,
		},
		{
			name: "pending task cannot increase actual pomodoros",
			setup: func() *Task {
				return NewTask("Write tests", 3)
			},
			wantActual: 0,
		},
		{
			name: "done task cannot increase actual pomodoros",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				task.StartWork()
				task.IncreaseActualPomos()
				task.MarkDone()
				return task
			},
			wantActual: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setup()
			task.IncreaseActualPomos()

			if task.ActualPomos != tt.wantActual {
				t.Errorf("ActualPomos = %d, want %d", task.ActualPomos, tt.wantActual)
			}
		})
	}
}

func TestTaskStopWork(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *Task
		wantStatus Status
	}{
		{
			name: "in-progress task can be stopped",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				task.StartWork()
				return task
			},
			wantStatus: Pending,
		},
		{
			name: "stopping a pending task is a no-op",
			setup: func() *Task {
				return NewTask("Write tests", 3)
			},
			wantStatus: Pending,
		},
		{
			name: "stopping a done task is a no-op",
			setup: func() *Task {
				task := NewTask("Write tests", 3)
				task.StartWork()
				task.MarkDone()
				return task
			},
			wantStatus: Done,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setup()
			task.StopWork()

			if task.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
			}
		})
	}
}

func TestTaskEdit(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() *Task
		editName     string
		editEstimate int
		wantName     string
		wantEstimate int
	}{
		{
			name: "edit updates name and estimate",
			setup: func() *Task {
				return NewTask("Old name", 3)
			},
			editName:     "New name",
			editEstimate: 5,
			wantName:     "New name",
			wantEstimate: 5,
		},
		{
			name: "edit with negative estimate defaults to 1",
			setup: func() *Task {
				return NewTask("Old name", 3)
			},
			editName:     "New name",
			editEstimate: -2,
			wantName:     "New name",
			wantEstimate: 1,
		},
		{
			name: "edit with zero estimate defaults to 1",
			setup: func() *Task {
				return NewTask("Old name", 3)
			},
			editName:     "New name",
			editEstimate: 0,
			wantName:     "New name",
			wantEstimate: 1,
		},
		{
			name: "edit with empty name keeps original name",
			setup: func() *Task {
				return NewTask("Old name", 3)
			},
			editName:     "",
			editEstimate: 5,
			wantName:     "Old name",
			wantEstimate: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.setup()
			task.Edit(tt.editName, tt.editEstimate)

			if task.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", task.Name, tt.wantName)
			}
			if task.EstimatedPomos != tt.wantEstimate {
				t.Errorf("EstimatedPomos = %d, want %d", task.EstimatedPomos, tt.wantEstimate)
			}
		})
	}
}
