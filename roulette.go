package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("currId: %v\n", m.currId))
	return s.String()
}
