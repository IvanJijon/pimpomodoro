package theme

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Theme struct {
	Work       string `yaml:"work"`
	ShortBreak string `yaml:"short-break"`
	LongBreak  string `yaml:"long-break"`
	Paused     string `yaml:"paused"`
	Subtle     string `yaml:"subtle"`
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

func LoadFromFile(path string) Theme {
	content, err := os.ReadFile(path)
	if err != nil {
		return DefaultTheme()
	}

	var theme Theme
	err = yaml.Unmarshal(content, &theme)
	if err != nil {
		return DefaultTheme()
	}

	theme.fillMissingColors()

	return theme
}

func (th *Theme) fillMissingColors() {
	defaultTheme := DefaultTheme()

	if th.Work == "" {
		th.Work = defaultTheme.Work
	}
	if th.ShortBreak == "" {
		th.ShortBreak = defaultTheme.ShortBreak
	}
	if th.LongBreak == "" {
		th.LongBreak = defaultTheme.LongBreak
	}
	if th.Paused == "" {
		th.Paused = defaultTheme.Paused
	}
	if th.Subtle == "" {
		th.Subtle = defaultTheme.Subtle
	}
}
