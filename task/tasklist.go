package task

import (
	"sort"
)

type TaskList struct {
	tasks []*Task
}

func NewTaskList() *TaskList {
	return &TaskList{
		tasks: []*Task{},
	}
}

func (tl *TaskList) Tasks() []*Task {
	return tl.tasks
}

func (tl *TaskList) Add(task *Task) {
	tl.tasks = append(tl.tasks, task)
}

func (tl *TaskList) Remove(task *Task) {
	for i, t := range tl.tasks {
		if t == task {
			tl.tasks = append(tl.tasks[:i], tl.tasks[i+1:]...)
			break
		}
	}
}

func (tl *TaskList) Len() int {
	return len(tl.tasks)
}

func (tl *TaskList) SelectWIP(t *Task) {
	changed := false
	for _, task := range tl.tasks {
		if task == t {
			task.StartWork()
			changed = true
			break
		}
	}

	if changed {
		for _, task := range tl.tasks {
			if task != t && task.Status == InProgress {
				task.StopWork()
			}
		}
	}

	tl.sort()
}

func (tl *TaskList) CurrentWIP() *Task {
	for _, task := range tl.tasks {
		if task.Status == InProgress {
			return task
		}
	}

	return nil
}

func (tl *TaskList) MarkTaskDone(t *Task) {
	for _, task := range tl.tasks {
		if task == t {
			task.MarkDone()
			break
		}
	}

	tl.sort()
}

// statusPriority defines the order of task statuses for sorting: InProgress first, then Pending, and Done last.
var statusPriority = map[Status]int{
	InProgress: 0,
	Pending:    1,
	Done:       2,
}

func (tl *TaskList) sort() {
	sort.Slice(tl.tasks, func(i, j int) bool {
		return statusPriority[tl.tasks[i].Status] < statusPriority[tl.tasks[j].Status]
	})
}
