package machine

import (
	"fmt"
	"strings"
)

type Machine struct {
	Lights  string
	Buttons [][]int
	Joltage []int
}

func NewMachine(data string) Machine {
	parts := strings.Split(data, " ")

	machine := &Machine{
		Lights:  strings.Trim(parts[0], "[]"),
		Buttons: [][]int{},
		Joltage: []int{},
	}

	for _, part := range parts[1:] {
		if strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") {
			var buttonIndices []int
			buttonStr := strings.Trim(part, "()")
			for _, b := range strings.Split(buttonStr, ",") {
				var index int
				fmt.Sscanf(b, "%d", &index)
				buttonIndices = append(buttonIndices, index)
			}
			machine.Buttons = append(machine.Buttons, buttonIndices)
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			var jolt int
			joltStr := strings.Trim(part, "{}")
			for _, j := range strings.Split(joltStr, ",") {
				fmt.Sscanf(j, "%d", &jolt)
				machine.Joltage = append(machine.Joltage, jolt)
			}
		}
	}
	return *machine
}

func (m *Machine) IsOn() bool {
	return strings.Count(m.Lights, "#") == 0
}

func (m *Machine) IsPowered() bool {
	for _, j := range m.Joltage {
		if j != 0 {
			return false
		}
	}
	return true
}

func (m *Machine) Toggle(buttonIndex int) {
	if buttonIndex < 0 || buttonIndex >= len(m.Buttons) {
		return
	}
	button := m.Buttons[buttonIndex]

	// Toggle lights
	runes := []rune(m.Lights)
	for _, idx := range button {
		if idx < 0 || idx >= len(runes) {
			continue
		}
		if runes[idx] == '#' {
			runes[idx] = '.'
		} else {
			runes[idx] = '#'
		}
	}
	m.Lights = string(runes)
}

func (m *Machine) ToggleMask(joltageMultiplier []int) {
	for i := range m.Joltage {
		if i < len(joltageMultiplier) {
			m.Joltage[i] = (m.Joltage[i] - joltageMultiplier[i]) / 2
		}
	}
}
