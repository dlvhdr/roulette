package app

import (
	"fmt"
	"math"
)

func (m Model) printDebugInfo() string {
	return fmt.Sprintf(
		"%v\nPos: %v, Vel: %v\nRoll: %v, Equalibrium Pos: %v\nResIdx: %v\nCurrIdx: %v, Current Value: %v\nMiddle: %v\nWinner: %v\n\n",
		chars[m.projectToProgressBar()],
		m.pos,
		m.velocity,
		m.roll,
		m.res,
		m.res%len(m.options),
		int(math.Min(float64(len(m.options))-1, float64(len(m.options)-(m.roll+len(m.options)/2)%len(m.options)))),
		m.options[int(math.Min(float64(len(m.options))-1, float64(len(m.options)-(m.roll+len(m.options)/2)%len(m.options))))],
		len(m.options)/2,
		m.options[m.res%len(m.options)],
	)
}

func (m Model) projectToProgressBar() int {
	input_start := 0.0
	input_end := float64(m.res)
	output_start := 0.0
	output_end := math.Min(
		math.Max(
			0.0,
			math.RoundToEven(float64(len(chars)/2))+float64(m.res)-m.pos,
		),
		float64(len(chars)-1),
	)
	slope := (output_end - output_start) / (input_end - input_start)
	output := output_start + slope*(m.pos-input_start)
	return int(math.RoundToEven(output))
}
