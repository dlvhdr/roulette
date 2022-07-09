package app

import (
	"fmt"
	"math"
)

func (m Model) printDebugInfo() string {
	middle := (m.height / 2)
	currPos := (middle + m.currRoll) % m.height
	finalPos := (middle + m.totalRolls) % m.height
	itemsPassed := (finalPos / itemHeight)
	winner := m.options[len(m.options)-1-itemsPassed]
	return fmt.Sprintf(
		"%v\nPosition: %v, Velocity: %v\ncurrRoll: %v, totalRolls: %v\nwinnerIdx: %v, winner: %v\n\nheight: %v,finalPos: %v, middle: %v\nitemsPassed: %v,currPos: %v\n",
		chars[m.projectToProgressBar()],
		m.position,
		m.velocity,
		m.currRoll,
		m.totalRolls,
		m.winnerIdx,
		winner,
		m.height,
		finalPos,
		middle,
		itemsPassed,
		currPos,
	)
}

func (m Model) projectToProgressBar() int {
	input_start := 0.0
	input_end := float64(m.totalRolls)
	output_start := 0.0
	output_end := math.Min(
		math.Max(
			0.0,
			math.RoundToEven(float64(len(chars)/2))+float64(m.totalRolls)-m.position,
		),
		float64(len(chars)-1),
	)
	slope := (output_end - output_start) / (input_end - input_start)
	output := output_start + slope*(m.position-input_start)
	return int(math.RoundToEven(output))
}
