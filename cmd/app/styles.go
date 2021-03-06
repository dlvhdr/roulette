package app

import "github.com/charmbracelet/lipgloss"

const (
	gameWidth  = 30
	gameHeight = 13
	itemHeight = 3
	itemWidth  = 20
)

var (
	optionContainer = lipgloss.NewStyle().
			Width(itemWidth).
			Height(itemHeight).
			Padding(1, 0).
			Foreground(lipgloss.Color("0")).
			BorderForeground(lipgloss.Color("0")).
			Align(lipgloss.Center)

	titleContainer = lipgloss.NewStyle().Bold(true).Padding(1, 0)

	pointerContainer = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	colors = []string{
		"#3d59a1",
		"#f7768e",
		"#0db9d7",
		"#ff9e64",
		"#B4F9F8",
		"#9d7cd8",
		"#394b70",
		"#bb9af7",
		"#ff007c",
		"#7dcfff",
		"#2ac3de",
		"#e0af68",
		"#89ddff",
		"#9ece6a",
		"#73daca",
		"#41a6b5",
		"#7aa2f7",
		"#1abc9c",
	}

	chars = []string{"██", "▇▇", "▆▆", "▅▅", "▄▄", "▃▃", "▂▂", "▁▁", "__", "▁▁", "▂▂", "▃▃", "▄▄", "▅▅", "▆▆", "▇▇", "██"}
)
