package theme

type Theme struct {
	Work       string
	ShortBreak string
	LongBreak  string
	Paused     string
	Subtle     string
}

func DefaultTheme() Theme {
	return Theme{
		Work:       "#722F37",
		ShortBreak: "#40E0D0",
		LongBreak:  "#1E3A5F",
		Paused:     "#FFD700",
		Subtle:     "#666666",
	}
}
