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
	fps              = 60
	angularFrequency = 100.0
	damping          = 50.0
	stepsMultiplier  = 50
)

type Model struct {
	title          string
	options        []string
	currStep       int
	totalSteps     int
	spring         harmonica.Spring
	velocity       float64
	position       float64
	lightsOn       bool
	middlePosition int
	winnerIdx      int
	showWinner     bool
	help           help.Model
	debug          bool
}

func InitialModel(title string, options []string, debug bool) Model {
	rand.Seed(time.Now().UnixNano())
	numOptions := len(options)
	winnerIdx := rand.Intn(numOptions)

	optionsHeight := numOptions * itemHeight
	minSteps := optionsHeight * stepsMultiplier

	var middle int
	if optionsHeight%2 != 0 {
		middle = int(math.Floor(float64(optionsHeight) / 2.0))
	} else {
		middle = optionsHeight / 2
	}

	distanceToWinner := middle - winnerIdx*itemHeight
	if winnerIdx >= middle {
		distanceToWinner += optionsHeight
	}
	totalSteps := minSteps + distanceToWinner

	return Model{
		title:          title,
		options:        options,
		totalSteps:     totalSteps,
		debug:          debug,
		middlePosition: middle,
		winnerIdx:      winnerIdx,
		spring:         harmonica.NewSpring(harmonica.FPS(fps), angularFrequency, damping),
		help:           help.NewModel(),
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
			if m.currStep == 0 {
				return m, tea.Batch(m.flicker(), m.doRoll())
			} else {
				m.currStep = m.totalSteps
				m.showWinner = true
				return m, tea.Quit
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case stopGame:
		m.showWinner = true
		return m, tea.Quit

	case rollMsg:
		newPos, newVel := m.spring.Update(m.position, m.velocity, float64(m.totalSteps))
		m.velocity = newVel
		m.currStep = int(math.Round(newPos))
		m.position = newPos
		if floatEquals(newPos, float64(m.totalSteps)) {
			m.lightsOn = true
			m.showWinner = true
			return m, tea.Batch(m.flashWinner(), m.quit())
		}
		return m, m.doRoll()

	case flickerLightsMsg:
		m.lightsOn = !m.lightsOn
		return m, m.flicker()

	case flashWinner:
		m.showWinner = !m.showWinner
		return m, m.flashWinner()

	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(m.renderTitle())
	s.WriteString("\n")

	options := m.renderOptions()
	lights := m.renderLights()
	game := lipgloss.PlaceHorizontal(gameWidth, lipgloss.Center, lipgloss.JoinHorizontal(
		lipgloss.Center,
		lights,
		pointerContainer.Render(" ▶ "),
		options,
		pointerContainer.Render("  "),
		lights,
	))
	s.WriteString(game)
	s.WriteString("\n\n")

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
		return flickerLightsMsg{}
	})
}

type flickerLightsMsg struct{}

type rollMsg time.Time

func (m Model) renderLights() string {
	lights := make([]string, gameHeight)
	for i := 0; i < gameHeight; i++ {
		var light string
		if m.lightsOn {
			light = "✦"
		} else {
			light = "✧"
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
		gameWidth,
		lipgloss.Center,
		titleContainer.Render(m.title),
	)
}

func (m Model) renderOptions() string {
	renderedOptions := make([]string, len(m.options))
	for i, option := range m.options {
		color := lipgloss.Color(colors[i])
		container := optionContainer.Copy().Background(color).BorderBackground(color)
		if i == m.winnerIdx && m.showWinner {
			container = container.
				Border(lipgloss.ThickBorder()).
				Padding(0, 0).
				Height(1).
				Width(itemWidth - 2)
		}
		renderedOptions[i] = container.Render(fmt.Sprintf("%s", option))
	}
	optionsRows := strings.Split(
		lipgloss.JoinVertical(
			lipgloss.Center,
			renderedOptions...,
		),
		"\n",
	)

	optionsViewport := make([]string, len(optionsRows))
	for i := range optionsRows {
		optionsViewport[(i+m.currStep)%len(optionsRows)] = optionsRows[i]
	}

	halfGameHeight := gameHeight / 2
	low := Max(0, m.middlePosition-halfGameHeight)
	high := Min(len(optionsViewport), m.middlePosition+halfGameHeight+1)
	optionsViewport = optionsViewport[low:high]
	return lipgloss.PlaceVertical(
		gameHeight,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, optionsViewport...),
	)
}

type keyMap struct {
	Roll key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Roll: key.NewBinding(
		key.WithKeys("Enter"),
		key.WithHelp("Enter", "Roll/Speed Up"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "Quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Roll, keys.Quit}
}

func (m Model) quit() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		return stopGame{}
	}
}

func (m Model) flashWinner() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return flashWinner{}
	})
}

type flashWinner struct{}

type stopGame struct{}
