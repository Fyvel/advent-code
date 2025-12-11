package machine

import (
	"fmt"
	"strings"
)

type Machine struct {
	Lights     string
	Buttons    [][]int
	JoltageMap map[int]bool
	IsOn       bool
}

func NewMachine(data string) Machine {
	parts := strings.Split(data, " ")

	machine := &Machine{
		Lights:     strings.Trim(parts[0], "[]"),
		Buttons:    [][]int{},
		JoltageMap: make(map[int]bool),
		IsOn:       strings.Count(parts[0], "#") == 0,
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
				machine.JoltageMap[jolt] = true
			}
		}
	}
	return *machine
}

func (m *Machine) isLightOn() bool {
	return strings.Count(m.Lights, "#") == 0
}

func (m *Machine) Toggle(buttonIndex int) {
	if buttonIndex < 0 || buttonIndex >= len(m.Buttons) {
		return
	}
	button := m.Buttons[buttonIndex]
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
	m.IsOn = m.isLightOn()
}
