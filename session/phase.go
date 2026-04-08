package session

type Phase int

const (
	Idle Phase = iota
	Work
	ShortBreak
	LongBreak
)
