package task

type Status int

const (
	Pending Status = iota
	InProgress
	Done
)

type Task struct {
	Name           string
	EstimatedPomos int
	ActualPomos    int
	Status         Status
}

func NewTask(name string, estimatedPomos int) *Task {
	return &Task{
		Name:           name,
		EstimatedPomos: estimatedPomos,
		ActualPomos:    0,
		Status:         Pending,
	}
}

func (t *Task) StartWork() {
	if t.Status == Pending {
		t.Status = InProgress
	}
}

func (t *Task) StopWork() {
	if t.Status == InProgress {
		t.Status = Pending
	}
}

func (t *Task) MarkDone() {
	if t.Status == InProgress {
		t.Status = Done
	}
}

func (t *Task) IncreaseActualPomos() {
	if t.Status == InProgress {
		t.ActualPomos++
	}
}
