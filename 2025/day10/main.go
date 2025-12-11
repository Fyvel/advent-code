package main

import (
	"aoc2025/day10/machine"
	"aoc2025/day10/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatData(rows []string) []string {
	return rows
}

func interactive(data []string, v *utils.Visualiser) {
	m := machine.NewMachine(data[0])
	v.Render(m.Lights, m.IsOn, m.Buttons, -1)
	for !m.IsOn {
		var input int
		fmt.Print("Press button index: ")
		_, err := fmt.Scanf("%d\n", &input)
		if err != nil {
			fmt.Println("lol")
			break
		}
		m.Toggle(input)
		v.Render(m.Lights, m.IsOn, m.Buttons, input)
	}
}

type State struct {
	lights   string
	sequence []int
	pressed  uint
}

func findBestSequenceBFS(line string, v *utils.Visualiser, machineIdx int) (int, []int) {
	m := machine.NewMachine(line)
	if m.IsOn {
		return 0, []int{}
	}

	queue := []State{{m.Lights, []int{}, 0}}
	visited := make(map[string]bool)
	visited[m.Lights] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for b := range len(m.Buttons) {
			// check if already pressed
			if (current.pressed & (1 << b)) != 0 {
				continue
			}

			testMachine := machine.NewMachine(line)
			testMachine.Lights = current.lights
			testMachine.Toggle(b)

			if v != nil && machineIdx >= 0 {
				v.Update(machineIdx, testMachine.Lights, testMachine.IsOn, testMachine.Buttons, b)
			}

			if testMachine.IsOn {
				return len(current.sequence) + 1, append(current.sequence, b)
			}

			if !visited[testMachine.Lights] {
				visited[testMachine.Lights] = true
				queue = append(queue, State{
					testMachine.Lights,
					append(append([]int{}, current.sequence...), b),
					current.pressed | (1 << b),
				})
			}
		}
	}

	return -1, nil
}

type ActiveMachine struct {
	idx  int
	line string
}

func part1(data []string, v *utils.Visualiser) {
	if v == nil {
		return
	}

	var activeMachines []ActiveMachine
	for idx, line := range data {
		// map machine/line
		v.RegisterMachine(idx)
		activeMachines = append(activeMachines, ActiveMachine{idx, line})
	}

	// placeholders
	for range activeMachines {
		fmt.Println()
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	totalCount := 0
	results := make(map[string][]int)

	for _, vm := range activeMachines {
		wg.Add(1)
		go func(machineIdx int, line string) {
			defer wg.Done()

			pressCount, sequence := findBestSequenceBFS(line, v, machineIdx)

			mu.Lock()
			results[line] = sequence
			if pressCount != -1 {
				totalCount += pressCount
			}
			mu.Unlock()

			// mark complete
			m := machine.NewMachine(line)
			v.Complete(machineIdx, m.Lights, m.Buttons)
		}(vm.idx, vm.line)
	}

	wg.Wait()

	// RESULTS OUTPUT
	showResults(results, totalCount)

	fmt.Println("Part 1:", totalCount)
}

func part2(data []string) {
	// fmt.Println("Part 2:", data)
}

func showResults(results map[string][]int, totalPresses int) {
	// longest string to pad right
	maxMachineLen := len("Machines")
	for machineData := range results {
		start := strings.Index(machineData, "[")
		end := strings.Index(machineData, "]")
		if start != -1 && end != -1 && end > start {
			lightData := machineData[start+1 : end]
			if len(lightData) > maxMachineLen {
				maxMachineLen = len(lightData)
			}
		}
	}

	// headers
	fmt.Printf("\n┌─%s─┬─────────────┐\n", strings.Repeat("─", maxMachineLen))
	fmt.Printf("│ %-*s │ Press Count │\n", maxMachineLen, "Machines")
	fmt.Printf("├─%s─┼─────────────┤\n", strings.Repeat("─", maxMachineLen))
	// contents
	for machineData, sequence := range results {
		start := strings.Index(machineData, "[")
		end := strings.Index(machineData, "]")
		lightData := ""
		if start != -1 && end != -1 && end > start {
			lightData = machineData[start+1 : end]
		}
		fmt.Printf("│ %-*s │ %-11d │\n", maxMachineLen, lightData, len(sequence))
	}
	fmt.Printf("├─%s─┼─────────────┤\n", strings.Repeat("─", maxMachineLen))
	// footer
	fmt.Printf("│ %-*s │ %-11d │\n", maxMachineLen, "Total", totalPresses)
	fmt.Printf("└─%s─┴─────────────┘\n", strings.Repeat("─", maxMachineLen))
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	renderer := utils.NewVisualiser(1*time.Millisecond, false)

	interactive(formattedData, renderer)
	part1(formattedData, renderer)
	part2(formattedData)
}
