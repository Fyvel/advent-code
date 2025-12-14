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

func part1(data []string, v *utils.Visualiser) {
	if v == nil {
		return
	}

	for idx := range data {
		// map machine/line
		v.RegisterMachine(idx)
		// placeholders
		fmt.Println()
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	totalCount := 0
	results := make(map[string][]int)

	for idx, line := range data {
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
		}(idx, line)
	}

	wg.Wait()

	// RESULTS OUTPUT
	showResults(results, totalCount, false)
	fmt.Println("Part 1:", totalCount)
}

func interactivePart2(data []string, v *utils.Visualiser) {
	m := machine.NewMachine(data[0])

	renderJoltageString := func(joltage []int) string {
		parts := make([]string, len(joltage))
		for i, val := range joltage {
			parts[i] = fmt.Sprintf("%d", val)
		}
		return fmt.Sprintf("{%s}", strings.Join(parts, ","))
	}

	v.Render(renderJoltageString(m.Joltage), m.IsPowered(), m.Buttons, -1)

	for !m.IsPowered() {
		var input int
		fmt.Print("Press button index: ")
		if _, err := fmt.Scanf("%d\n", &input); err != nil {
			fmt.Println("lol")
			break
		}
		m.Toggle(input)
		v.Render(renderJoltageString(m.Joltage), m.IsPowered(), m.Buttons, input)
	}
}

func part2(data []string, v *utils.Visualiser) {
	if v == nil {
		return
	}

	for idx := range data {
		// map machine/line
		v.RegisterMachine(idx)
		// placeholders
		fmt.Println()
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	totalCount := 0
	results := make(map[string][]int)

	for idx, line := range data {
		wg.Add(1)
		go func(machineIdx int, line string) {
			defer wg.Done()
			m := machine.NewMachine(line)

			var pressCount int
			var sequence []int

			if m.IsPowered() {
				pressCount = 0
				sequence = []int{}
				mu.Lock()
				results[line] = sequence
				mu.Unlock()
				v.CompleteJoltage(machineIdx, m.Joltage, m.Buttons)
				return
			}

			patterns := buildButtonConfig(m)
			memo := make(map[string]*joltageResult)

			if result := solveJoltageDFS(m.Joltage, patterns, memo, v, machineIdx, m.Buttons); result != nil {
				pressCount = result.presses
				sequence = result.sequence

			}

			mu.Lock()
			results[line] = sequence
			if pressCount != -1 {
				totalCount += pressCount
			}
			mu.Unlock()

			// mark complete
			v.CompleteJoltage(machineIdx, m.Joltage, m.Buttons)
		}(idx, line)
	}

	wg.Wait()

	// RESULTS OUTPUT
	showResults(results, totalCount, true)

	fmt.Println("Part 2:", totalCount)
}

func showResults(results map[string][]int, totalPresses int, isPart2 bool) {
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
	if !isPart2 {
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
	}
	// footer
	fmt.Printf("│ %-*s │ %-11d │\n", maxMachineLen, "Total", totalPresses)
	fmt.Printf("└─%s─┴─────────────┘\n", strings.Repeat("─", maxMachineLen))
	time.Sleep(1 * time.Second)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	renderer := utils.NewVisualiser(1*time.Millisecond, false)

	interactivePart1(formattedData, renderer)
	part1(formattedData, renderer)
	interactivePart2(formattedData, renderer)
	part2(formattedData, renderer)
}

// #region Part 1

func interactivePart1(data []string, v *utils.Visualiser) {
	m := machine.NewMachine(data[0])
	v.Render(m.Lights, m.IsOn(), m.Buttons, -1)
	for !m.IsOn() {
		var input int
		fmt.Print("Press button index: ")
		_, err := fmt.Scanf("%d\n", &input)
		if err != nil {
			fmt.Println("lol")
			break
		}
		m.Toggle(input)
		v.Render(m.Lights, m.IsOn(), m.Buttons, input)
	}
}

type lightsResult struct {
	lights   string
	sequence []int
	pressed  uint
}

func findBestSequenceBFS(line string, v *utils.Visualiser, machineIdx int) (int, []int) {
	m := machine.NewMachine(line)
	if m.IsOn() {
		return 0, []int{}
	}

	queue := []lightsResult{{m.Lights, []int{}, 0}}
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
				v.Update(machineIdx, testMachine.Lights, testMachine.IsOn(), testMachine.Buttons, b)
			}

			if testMachine.IsOn() {
				return len(current.sequence) + 1, append(current.sequence, b)
			}

			if !visited[testMachine.Lights] {
				visited[testMachine.Lights] = true
				queue = append(queue, lightsResult{
					testMachine.Lights,
					append(current.sequence, b),
					current.pressed | (1 << b),
				})
			}
		}
	}

	return -1, nil
}

// #endregion Part 1

// #region Part 2
type joltageResult struct {
	presses  int
	sequence []int
}

type joltageButtonConfig struct {
	joltageMultiplier []int
	presses           int
	buttonMask        int
}

func buildButtonConfig(m machine.Machine) map[string][]joltageButtonConfig {
	patterns := make(map[string][]joltageButtonConfig)
	numLights := len(m.Joltage)

	// all button combinations (2^n possibilities)
	for n := 0; n < (1 << len(m.Buttons)); n++ {
		lightResult := make([]int, numLights)
		joltageMultiplier := make([]int, numLights)
		presses := 0

		for buttonIndex := 0; buttonIndex < len(m.Buttons); buttonIndex++ {
			if (n & (1 << buttonIndex)) != 0 {
				btn := m.Buttons[buttonIndex]
				for _, light := range btn {
					if light < numLights {
						lightResult[light] ^= 1
						joltageMultiplier[light]++
					}
				}
				presses++
			}
		}

		key := fmt.Sprintf("%v", lightResult)
		patterns[key] = append(patterns[key], joltageButtonConfig{
			joltageMultiplier: joltageMultiplier,
			presses:           presses,
			buttonMask:        n,
		})
	}

	return patterns
}

func solveJoltageDFS(joltages []int, buttonConfigMap map[string][]joltageButtonConfig, memo map[string]*joltageResult, v *utils.Visualiser, machineIdx int, buttons [][]int) *joltageResult {
	// base case
	if isPowered(joltages) {
		if v != nil && machineIdx >= 0 {
			v.UpdateJoltage(machineIdx, joltages, true, buttons, -1)
		}
		return &joltageResult{presses: 0, sequence: []int{}}
	}

	// check memo
	key := fmt.Sprintf("%v", joltages)
	if val, ok := memo[key]; ok {
		return val
	}

	// least significant bit from each joltage value
	lsb := make([]int, len(joltages))
	for i, j := range joltages {
		lsb[i] = j & 1
	}

	lsbKey := fmt.Sprintf("%v", lsb)
	if combinations, ok := buttonConfigMap[lsbKey]; ok {
		var bestResult *joltageResult

		for _, c := range combinations {
			testMachine := machine.Machine{
				Buttons: buttons,
				Joltage: make([]int, len(joltages)),
			}
			copy(testMachine.Joltage, joltages)
			testMachine.ToggleMask(c.joltageMultiplier)
			newJoltage := testMachine.Joltage

			if v != nil && machineIdx >= 0 {
				firstButton := -1
				for i := range 32 {
					if (c.buttonMask & (1 << i)) != 0 {
						firstButton = i
						break
					}
				}
				v.UpdateJoltage(machineIdx, newJoltage, testMachine.IsPowered(), buttons, firstButton)
			}

			// termination
			if !isValidJoltage(newJoltage) {
				continue
			}

			// recursion
			if rest := solveJoltageDFS(newJoltage, buttonConfigMap, memo, v, machineIdx, buttons); rest != nil {
				totalPresses := c.presses + 2*rest.presses

				currentButtons := maskToButton(c.buttonMask)
				newSequence := append(currentButtons, rest.sequence...)

				if bestResult == nil || totalPresses < bestResult.presses {
					bestResult = &joltageResult{
						presses:  totalPresses,
						sequence: newSequence,
					}
				}
			}
		}

		// store
		memo[key] = bestResult
		return bestResult
	}

	memo[key] = nil
	return nil
}

func isValidJoltage(joltages []int) bool {
	for _, j := range joltages {
		if j < 0 {
			return false
		}
	}
	return true
}

func isPowered(joltages []int) bool {
	for _, j := range joltages {
		if j != 0 {
			return false
		}
	}
	return true
}

func maskToButton(mask int) []int {
	var buttons []int
	for i := 0; mask > 0; i++ {
		if mask&1 == 1 {
			buttons = append(buttons, i)
		}
		mask >>= 1
	}
	return buttons
}

// #endregion Part 2
