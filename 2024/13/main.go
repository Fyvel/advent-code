package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

type Vector struct {
	x int
	y int
}
type Button struct {
	name   string
	tokens int
	move   Vector
}
type Machine struct {
	buttonA Button
	buttonB Button
	prize   Vector
}

func getLineCoords(line string) Vector {
	re := regexp.MustCompile(`X[+=](\d+), Y[+=](\d+)`)
	matches := re.FindStringSubmatch(line)
	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])
	return Vector{x: x, y: y}
}

func formatData(rows []string) []Machine {
	machines := []Machine{}
	for i := 0; i < len(rows); i += 4 {
		lineButtonA := strings.Split(rows[i], "Button A: ")[1]
		lineButtonB := strings.Split(rows[i+1], "Button B: ")[1]
		linePrize := strings.Split(rows[i+2], "Prize: ")[1]
		machine := Machine{
			buttonA: Button{
				name:   "A",
				tokens: 3,
				move:   getLineCoords(lineButtonA),
			},
			buttonB: Button{
				name:   "B",
				tokens: 1,
				move:   getLineCoords(lineButtonB),
			},
			prize: getLineCoords(linePrize),
		}
		machines = append(machines, machine)
	}
	return machines
}

func calculateMinimumTokens(buttons []Button, prize Vector, memo map[string]int) int {
	key := fmt.Sprintf("%d_%d", prize.x, prize.y)
	if value, ok := memo[key]; ok {
		return value
	}

	if prize.x == 0 && prize.y == 0 {
		return 0
	}

	if prize.x < 0 || prize.y < 0 {
		memo[key] = -1
		return -1
	}

	minTokens := math.MaxInt32

	for _, button := range buttons {
		newTarget := Vector{
			x: prize.x - button.move.x,
			y: prize.y - button.move.y,
		}
		tokenCount := calculateMinimumTokens(buttons, newTarget, memo)

		if tokenCount != -1 && tokenCount+1 < minTokens {
			minTokens = tokenCount + button.tokens
		}
	}

	if minTokens == math.MaxInt32 {
		memo[key] = -1
	} else {
		memo[key] = minTokens
	}

	return memo[key]
}

func part1(machines []Machine) int {
	tokens := 0

	for idx, machine := range machines {
		if idx%40 == 0 {
			fmt.Printf("Machine %d out of %d\n", idx+1, len(machines))
		}

		memo := make(map[string]int)
		minMachineTokens := calculateMinimumTokens([]Button{machine.buttonA, machine.buttonB}, machine.prize, memo)
		if minMachineTokens != -1 {
			tokens += minMachineTokens
		}
	}

	fmt.Println("Part 1:", tokens)
	return tokens
}

func calculateMinimumTokensMath(machine Machine) int {
	// equations
	// btnACount(ax) + btnBCount(bx) = px*constant
	// btnACount(ay) + btnBCount(by) = py*constant
	ax, ay := machine.buttonA.move.x, machine.buttonA.move.y
	bx, by := machine.buttonB.move.x, machine.buttonB.move.y
	px, py := machine.prize.x, machine.prize.y

	determinant := float64(ax*by - ay*bx)

	if determinant == 0 {
		return -1 // impossible to solve
	}

	// solve for btnACount and btnBCount
	// btnACount = (px*by - py*bx) / (ax*by - ay*bx)
	// btnBCount = (ax*py - ay*px) / (ax*by - ay*bx)
	btnACount := float64(px*by-py*bx) / determinant
	btnBCount := float64(ax*py-ay*px) / determinant

	// legit counts (must be positive and full numbers)
	if btnACount < 0 || btnBCount < 0 || math.Floor(btnACount) != btnACount || math.Floor(btnBCount) != btnBCount {
		return -1
	}

	// calculate the total tokens
	// tokens = btnACount * btnA.tokens + btnBCount * btnB.tokens
	return int(float64(machine.buttonA.tokens)*btnACount + float64(machine.buttonB.tokens)*btnBCount)
}

func part2(machines []Machine) int {
	tokens := 0

	for run := 0; run < 2; run++ {
		for idx, machine := range machines {
			if idx%40 == 0 {
				fmt.Printf("Machine %d out of %d\n", idx+1, len(machines))
			}

			constant := 0
			if run%2 != 0 {
				constant = 10000000000000
			}

			newPrizeX := constant + machine.prize.x
			newPrizeY := constant + machine.prize.y

			machine.prize = Vector{
				x: newPrizeX,
				y: newPrizeY,
			}

			minMachineTokens := calculateMinimumTokensMath(machine)
			if minMachineTokens != -1 {
				tokens += minMachineTokens
			}
		}

		fmt.Printf("Part %d: %d\n", run+1, tokens)
	}

	return tokens
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData)
	part2(formattedData)
}
