package app

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	help "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxHeight                = itemHeight * 5
	width                    = 30
	fps                      = 60
	epsilon          float64 = 0.5
	angularFrequency         = 100.0
	damping                  = 50.0
	minSpins                 = 10
)

type Model struct {
	title        string
	options      []string
	roll         int
	help         help.Model
	spring       harmonica.Spring
	velocity     float64
	pos          float64
	res          int
	dimmedLights bool
	debug        bool
}

func InitialModel(title string, options []string, debug bool) Model {
	rand.Seed(time.Now().UnixNano())
	res := minSpins*len(options) + rand.Intn(minSpins*len(options))
	return Model{
		title:        title,
		options:      options,
		roll:         0,
		help:         help.NewModel(),
		spring:       harmonica.NewSpring(harmonica.FPS(fps), angularFrequency, damping),
		velocity:     0.0,
		pos:          0.0,
		res:          res,
		debug:        debug,
		dimmedLights: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "enter":
			if m.roll == 0 {
				return m, tea.Batch(m.flicker(), m.doRoll())
			} else {
				return m, nil
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case rollMsg:
		newPos, newVel := m.spring.Update(m.pos, m.velocity, float64(m.res))
		m.velocity = newVel
		m.roll = int(math.Round(newPos))
		m.pos = newPos
		if floatEquals(newPos, float64(m.res)) {
			m.dimmedLights = false
			return m, tea.Quit
		}
		return m, m.doRoll()

	case flickerMsg:
		m.dimmedLights = !m.dimmedLights
		return m, m.flicker()
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(m.renderTitle())
	s.WriteString("\n")

	options := m.renderOptions()
	lights := m.renderLights(maxHeight)
	// put pointers in a 2 height container with alighment based on number of items
	game := lipgloss.NewStyle().Height(maxHeight).MaxHeight(maxHeight).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			lights,
			pointerContainer.Render(" ▶ "),
			options,
			pointerContainer.Render("  "),
			lights,
		),
	)
	s.WriteString(game)
	s.WriteString("\n\n")

	if m.debug {
		s.WriteString(m.printDebugInfo())
	}

	s.WriteString(m.help.ShortHelpView(keys.ShortHelp()))
	s.WriteString("\n")

	return s.String()
}

func (m Model) doRoll() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return rollMsg(t)
	})
}

func (m Model) flicker() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return flickerMsg{}
	})
}

type flickerMsg struct{}

type rollMsg time.Time

func (m Model) renderLights(height int) string {
	lights := make([]string, height)
	for i := 0; i < height; i++ {
		var light string
		if m.dimmedLights {
			light = "⚬"
		} else {
			light = "●"
		}
		lights[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0af68")).Render(light)
	}
	allLights := lipgloss.JoinVertical(lipgloss.Left, lights...)
	return allLights
}

func (m Model) renderTitle() string {
	if m.title == "" {
		return ""
	}

	return lipgloss.PlaceHorizontal(
		width,
		lipgloss.Center,
		titleContainer.Render(m.title),
	)
}

func (m Model) renderOptions() string {
	rendered := make([]string, len(m.options))
	for i := range m.options {
		option := m.options[i]
		newIdx := (i + m.roll) % len(m.options)
		color := lipgloss.Color(colors[i])
		container := optionContainer.Copy().Background(color).BorderBackground(color)
		rendered[newIdx] = container.Render(fmt.Sprintf("%s", option))
	}

	return lipgloss.JoinVertical(lipgloss.Center, rendered...)
}

type keyMap struct {
	Roll key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Roll: key.NewBinding(
		key.WithKeys("Enter"),
		key.WithHelp("Enter", "Roll"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "Stop the roll"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Roll, keys.Quit}
}

func floatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}
	return false
}
