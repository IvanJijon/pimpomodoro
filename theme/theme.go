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
	theme := DefaultTheme()

	content, err := os.ReadFile(path)
	if err != nil {
		return theme
	}

	if yaml.Unmarshal(content, &theme) != nil {
		return theme
	}

	return theme
}
