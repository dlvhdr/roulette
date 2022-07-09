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
	maxHeight                = 100000
	width                    = 30
	fps                      = 60
	epsilon          float64 = 0.5
	angularFrequency         = 100.0
	damping                  = 50.0
	itemHeight               = 3
	gameHeight               = 31
)

type Model struct {
	title        string
	options      []string
	currRoll     int
	totalRolls   int
	help         help.Model
	spring       harmonica.Spring
	velocity     float64
	position     float64
	dimmedLights bool
	debug        bool
	winnerIdx    int
	height       int
	middle       int
}

func InitialModel(title string, options []string, debug bool) Model {
	rand.Seed(time.Now().UnixNano())
	numOptions := len(options)
	winner := rand.Intn(numOptions)
	optionsHeight := numOptions * itemHeight
	minSpins := optionsHeight * 50
	fmt.Printf("Winner: %v\n", options[winner])

	var middle int
	fmt.Printf("Options height: %v\n", optionsHeight)
	if optionsHeight%2 != 0 {
		middle = int(math.Floor(float64(optionsHeight) / 2.0))
	} else {
		middle = optionsHeight / 2
	}
	fmt.Printf("middle: %v\n", middle)

	distanceToWinner := middle - winner*itemHeight
	if winner >= middle {
		distanceToWinner += optionsHeight
	}
	fmt.Printf("distanceToWinner: %v\n", distanceToWinner)
	totalRolls := minSpins + distanceToWinner

	// remainderRolls := minSpins % len(options)
	// distanceToWinner := int(math.Abs(float64(len(options) - 1 - winner - remainderRolls)))
	// totalRolls := minSpins + distanceToWinner

	return Model{
		title:        title,
		options:      options,
		currRoll:     0,
		help:         help.NewModel(),
		spring:       harmonica.NewSpring(harmonica.FPS(fps), angularFrequency, damping),
		velocity:     0.0,
		position:     0.0,
		totalRolls:   totalRolls,
		debug:        debug,
		dimmedLights: false,
		height:       13,
		middle:       middle,
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
			if m.currRoll == 0 {
				return m, tea.Batch(m.flicker(), m.doRoll())
			} else {
				return m, nil
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case rollMsg:
		newPos, newVel := m.spring.Update(m.position, m.velocity, float64(m.totalRolls))
		m.velocity = newVel
		m.currRoll = int(math.Round(newPos))
		m.position = newPos
		if floatEquals(newPos, float64(m.totalRolls)) {
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

	options := lipgloss.PlaceVertical(gameHeight, lipgloss.Center, m.renderOptions())
	lights := m.renderLights()
	game := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lights,
		pointerContainer.Render(" ▶ "),
		options,
		pointerContainer.Render("  "),
		lights,
	)
	s.WriteString(game)
	s.WriteString("\n\n")

	// if m.debug {
	// 	s.WriteString(m.printDebugInfo())
	// }

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

func (m Model) renderLights() string {
	lights := make([]string, gameHeight)
	for i := 0; i < gameHeight; i++ {
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
	for i, option := range m.options {
		color := lipgloss.Color(colors[i])
		container := optionContainer.Copy().Background(color).BorderBackground(color)
		// if option == m.calcWinner() {
		// 	container.Border(lipgloss.NormalBorder())
		// }
		rendered[i] = container.Render(fmt.Sprintf("%s", option))
	}
	joined := lipgloss.JoinVertical(lipgloss.Center, rendered...)
	split := strings.Split(joined, "\n")

	final := make([]string, len(split))
	for i := range split {
		final[(i+m.currRoll)%len(split)] = split[i]
	}

	numbers := make([]string, len(m.options)*itemHeight)
	for i := 0; i < len(m.options); i++ {
		for j := 0; j < itemHeight; j++ {
			idx := i*itemHeight + j
			if idx < m.middle {
				numbers[idx] = fmt.Sprintf("%d", idx-m.middle)
			} else if idx == m.middle {
				numbers[idx] = fmt.Sprint("0")
			} else {
				numbers[idx] = fmt.Sprintf("%d", idx-m.middle)
			}
		}
	}
	allNumbers := lipgloss.JoinVertical(lipgloss.Center, numbers...)
	allOptions := lipgloss.JoinVertical(lipgloss.Center, final...)

	return lipgloss.JoinHorizontal(lipgloss.Left, allNumbers, allOptions)
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

func (m Model) calcWinner() string {
	middle := (m.height / 2)
	finalPos := (middle + m.totalRolls) % m.height
	itemsPassed := (finalPos / itemHeight)
	winner := m.options[len(m.options)-1-itemsPassed]
	return winner
}
