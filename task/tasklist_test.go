package task

import "testing"

func TestTaskListAdd(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *TaskList
		wantCount int
	}{
		{
			name: "adding a task to an empty list",
			setup: func() *TaskList {
				tl := NewTaskList()
				tl.Add(NewTask("Write tests", 3))
				return tl
			},
			wantCount: 1,
		},
		{
			name: "adding multiple tasks",
			setup: func() *TaskList {
				tl := NewTaskList()
				tl.Add(NewTask("Write tests", 3))
				tl.Add(NewTask("Refactor model", 2))
				return tl
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := tt.setup()

			if tl.Len() != tt.wantCount {
				t.Errorf("Len() = %d, want %d", tl.Len(), tt.wantCount)
			}
		})
	}
}

func TestTaskListRemove(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (*TaskList, *Task)
		wantCount int
	}{
		{
			name: "removing a task decreases count",
			setup: func() (*TaskList, *Task) {
				tl := NewTaskList()
				task := NewTask("Write tests", 3)
				tl.Add(task)
				tl.Add(NewTask("Refactor model", 2))
				return tl, task
			},
			wantCount: 1,
		},
		{
			name: "removing a task not in the list is a no-op",
			setup: func() (*TaskList, *Task) {
				tl := NewTaskList()
				tl.Add(NewTask("Write tests", 3))
				outsider := NewTask("Not in list", 1)
				return tl, outsider
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl, task := tt.setup()
			tl.Remove(task)

			if tl.Len() != tt.wantCount {
				t.Errorf("Len() = %d, want %d", tl.Len(), tt.wantCount)
			}
		})
	}
}

func TestTaskListSelectWIP(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() (*TaskList, *Task)
		wantStatus   Status
		wantWIPCount int
	}{
		{
			name: "selecting a pending task sets it to in-progress",
			setup: func() (*TaskList, *Task) {
				tl := NewTaskList()
				task := NewTask("Write tests", 3)
				tl.Add(task)
				return tl, task
			},
			wantStatus:   InProgress,
			wantWIPCount: 1,
		},
		{
			name: "selecting a new task stops the previous WIP",
			setup: func() (*TaskList, *Task) {
				tl := NewTaskList()
				first := NewTask("Write tests", 3)
				second := NewTask("Refactor model", 2)
				tl.Add(first)
				tl.Add(second)
				tl.SelectWIP(first)
				return tl, second
			},
			wantStatus:   InProgress,
			wantWIPCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl, task := tt.setup()
			tl.SelectWIP(task)

			if task.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
			}

			wipCount := 0
			for _, tk := range tl.Tasks() {
				if tk.Status == InProgress {
					wipCount++
				}
			}
			if wipCount != tt.wantWIPCount {
				t.Errorf("WIP count = %d, want %d", wipCount, tt.wantWIPCount)
			}
		})
	}
}

func TestTaskListSelectWIPReorders(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() (*TaskList, *Task)
		wantFirstName string
	}{
		{
			name: "selected WIP task moves to top of list",
			setup: func() (*TaskList, *Task) {
				tl := NewTaskList()
				tl.Add(NewTask("First", 1))
				second := NewTask("Second", 2)
				tl.Add(second)
				tl.Add(NewTask("Third", 3))
				return tl, second
			},
			wantFirstName: "Second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl, task := tt.setup()
			tl.SelectWIP(task)

			first := tl.Tasks()[0]
			if first.Name != tt.wantFirstName {
				t.Errorf("first task = %q, want %q", first.Name, tt.wantFirstName)
			}
		})
	}
}

func TestTaskListMarkDoneReorders(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() *TaskList
		wantLastName string
	}{
		{
			name: "done task moves to bottom of list",
			setup: func() *TaskList {
				tl := NewTaskList()
				first := NewTask("First", 1)
				tl.Add(first)
				tl.Add(NewTask("Second", 2))
				tl.SelectWIP(first)
				tl.MarkTaskDone(first)
				return tl
			},
			wantLastName: "First",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := tt.setup()

			tasks := tl.Tasks()
			last := tasks[len(tasks)-1]
			if last.Name != tt.wantLastName {
				t.Errorf("last task = %q, want %q", last.Name, tt.wantLastName)
			}
		})
	}
}

func TestTaskListMarkDoneOutsider(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() (*TaskList, *Task)
		wantStatus Status
	}{
		{
			name: "marking done a task not in the list does not change it",
			setup: func() (*TaskList, *Task) {
				tl := NewTaskList()
				tl.Add(NewTask("In list", 1))
				outsider := NewTask("Not in list", 2)
				outsider.StartWork()
				tl.MarkTaskDone(outsider)
				return tl, outsider
			},
			wantStatus: InProgress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, outsider := tt.setup()

			if outsider.Status != tt.wantStatus {
				t.Errorf("outsider Status = %v, want %v", outsider.Status, tt.wantStatus)
			}
		})
	}
}

func TestTaskListCurrentWIP(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *TaskList
		wantNil  bool
		wantName string
	}{
		{
			name: "returns nil when no WIP task",
			setup: func() *TaskList {
				tl := NewTaskList()
				tl.Add(NewTask("Write tests", 3))
				return tl
			},
			wantNil: true,
		},
		{
			name: "returns the current WIP task",
			setup: func() *TaskList {
				tl := NewTaskList()
				task := NewTask("Write tests", 3)
				tl.Add(task)
				tl.SelectWIP(task)
				return tl
			},
			wantNil:  false,
			wantName: "Write tests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := tt.setup()
			wip := tl.CurrentWIP()

			if tt.wantNil && wip != nil {
				t.Errorf("CurrentWIP() = %v, want nil", wip)
			}
			if !tt.wantNil && wip == nil {
				t.Fatal("CurrentWIP() = nil, want a task")
			}
			if !tt.wantNil && wip.Name != tt.wantName {
				t.Errorf("CurrentWIP().Name = %q, want %q", wip.Name, tt.wantName)
			}
		})
	}
}
