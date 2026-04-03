package ui

type PimpomodoroUI struct {
	Header    *Header
	Countdown *Countdown
}

func NewPimpomodoroUI() PimpomodoroUI {
	return PimpomodoroUI{
		Header:    InitHeader(),
		Countdown: InitCountdown(),
	}
}
