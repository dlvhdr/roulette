package app

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) printDebugInfo() string {
	currPos := (m.middlePosition + m.currStep) % gameHeight
	finalPos := (m.middlePosition + m.totalSteps) % gameHeight
	itemsPassed := (finalPos / itemHeight)
	winner := m.options[len(m.options)-1-itemsPassed]
	return fmt.Sprintf(
		"%v\nPosition: %v, Velocity: %v\ncurrRoll: %v, totalRolls: %v\nwinnerIdx: %v, winner: %v\n\nheight: %v,finalPos: %v, middle: %v\nitemsPassed: %v,currPos: %v\n",
		chars[m.projectToProgressBar()],
		m.position,
		m.velocity,
		m.currStep,
		m.totalSteps,
		m.winnerIdx,
		winner,
		gameHeight,
		finalPos,
		m.middlePosition,
		itemsPassed,
		currPos,
	)
}

func (m Model) projectToProgressBar() int {
	input_start := 0.0
	input_end := float64(m.totalSteps)
	output_start := 0.0
	output_end := math.Min(
		math.Max(
			0.0,
			math.RoundToEven(float64(len(chars)/2))+float64(m.totalSteps)-m.position,
		),
		float64(len(chars)-1),
	)
	slope := (output_end - output_start) / (input_end - input_start)
	output := output_start + slope*(m.position-input_start)
	return int(math.RoundToEven(output))
}

func (m Model) printLineNumbers() string {

	numbers := make([]string, len(m.options)*itemHeight)
	for i := 0; i < len(m.options); i++ {
		for j := 0; j < itemHeight; j++ {
			idx := i*itemHeight + j
			if idx < m.middlePosition {
				numbers[idx] = fmt.Sprintf("%d", idx-m.middlePosition)
			} else if idx == m.middlePosition {
				numbers[idx] = fmt.Sprint("0")
			} else {
				numbers[idx] = fmt.Sprintf("%d", idx-m.middlePosition)
			}
		}
	}
	halfGameHeight := gameHeight / 2
	low := Max(0, m.middlePosition-halfGameHeight)
	high := Min(len(numbers), m.middlePosition+halfGameHeight+1)
	numbers = numbers[low:high]
	return lipgloss.JoinVertical(lipgloss.Center, numbers...)
}
