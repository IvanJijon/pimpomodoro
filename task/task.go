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
	if name == "" {
		name = "Untitled"
	}
	if estimatedPomos <= 0 {
		estimatedPomos = 1
	}

	return &Task{
		Name:           name,
		EstimatedPomos: estimatedPomos,
		ActualPomos:    0,
		Status:         Pending,
	}
}

func (t *Task) Edit(name string, estimatedPomos int) *Task {
	if name == "" {
		return t
	}
	if estimatedPomos <= 0 {
		estimatedPomos = 1
	}
	t.Name = name
	t.EstimatedPomos = estimatedPomos
	return t
}

func (t *Task) StartWork() {
	t.Status = InProgress
}

func (t *Task) StopWork() {
	if t.Status == InProgress {
		t.Status = Pending
	}
}

func (t *Task) MarkDone() {
	t.Status = Done
}

func (t *Task) UnmarkDone() {
	if t.Status == Done {
		t.Status = Pending
	}
}

func (t *Task) IncreaseActualPomos() {
	if t.Status == InProgress {
		t.ActualPomos++
	}
}
