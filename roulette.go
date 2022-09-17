package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	options []string
	currId  int
}

func InitialModel(options []string) Model {
	return Model{
		options: options,
		currId:  len(options) / 2,
	}
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case tickMsg:
		m.currId = (m.currId + 1) % len(m.options)
		return m, doTick()

	}

	return m, nil
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

var (
	itemStyle = lipgloss.NewStyle().
			Padding(0, 5).
			Width(20).
			Height(1).
			Align(lipgloss.Center)

	selectedItemStyle = itemStyle.Copy().
				Background(lipgloss.Color("1")).
				Foreground(lipgloss.Color("0")).
				Bold(true)
)

func (m Model) View() string {
	s := strings.Builder{}

	renderedOptions := make([]string, len(m.options))
	for i := 0; i < len(m.options); i++ {
		option := m.options[(i+m.currId)%len(m.options)]

		var style lipgloss.Style
		if i == len(m.options)/2 {
			style = selectedItemStyle
		} else {
			style = itemStyle
		}

		renderedOptions[i] = style.Render(option)
	}

	s.WriteString(lipgloss.JoinVertical(lipgloss.Center, renderedOptions...))

	s.WriteString("\n")
	return s.String()
}
