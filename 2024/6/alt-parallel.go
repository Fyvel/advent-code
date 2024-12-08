package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Position [2]int
type Direction string

var positionsMapping = map[Direction][2]int{
	"^": {-1, 0},
	"<": {0, -1},
	">": {0, 1},
	"v": {1, 0},
}

var moveMapping = map[Direction]Direction{
	"^": ">",
	"<": "^",
	">": "v",
	"v": "<",
}

type GameState struct {
	obstacles      map[string]bool
	areaLimits     map[string]bool
	guardPosition  Position
	guardDirection Direction
}

func initialise(grid [][]string) GameState {
	obstacles := make(map[string]bool)
	areaLimits := make(map[string]bool)
	var guardPosition Position
	var guardDirection Direction

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			cell := grid[i][j]
			key := fmt.Sprintf("%d_%d", i, j)

			if cell == "#" {
				obstacles[key] = true
			} else if i == len(grid)-1 || j == len(grid[0])-1 || i == 0 || j == 0 {
				areaLimits[key] = true
			}

			if cell == "^" || cell == "<" || cell == ">" || cell == "v" {
				guardPosition = Position{i, j}
				guardDirection = Direction(cell)
			}
		}
	}

	return GameState{
		obstacles:      obstacles,
		areaLimits:     areaLimits,
		guardPosition:  guardPosition,
		guardDirection: guardDirection,
	}
}

func simulate(start Position, direction Direction, newObstacle string, areaLimits, obstacles map[string]bool) bool {
	localMoves := make(map[string]map[Direction]bool)
	pos := start
	dir := direction

	for {
		key := fmt.Sprintf("%d_%d", pos[0], pos[1])
		if localMoves[key] != nil {
			if localMoves[key][dir] {
				return true
			}
		} else {
			localMoves[key] = make(map[Direction]bool)
		}
		localMoves[key][dir] = true

		nextPos := Position{
			pos[0] + positionsMapping[dir][0],
			pos[1] + positionsMapping[dir][1],
		}
		nextKey := fmt.Sprintf("%d_%d", nextPos[0], nextPos[1])

		if areaLimits[nextKey] && nextKey != newObstacle {
			return false
		}

		if obstacles[nextKey] || nextKey == newObstacle {
			dir = moveMapping[dir]
		} else {
			pos = nextPos
		}
	}
}

func part1(grid [][]string) int {
	state := initialise(grid)
	visited := make(map[string]bool)

	for {
		key := fmt.Sprintf("%d_%d", state.guardPosition[0], state.guardPosition[1])
		if state.areaLimits[key] {
			visited[key] = true
			break
		}

		visited[key] = true
		nextPosition := Position{
			state.guardPosition[0] + positionsMapping[state.guardDirection][0],
			state.guardPosition[1] + positionsMapping[state.guardDirection][1],
		}

		nextKey := fmt.Sprintf("%d_%d", nextPosition[0], nextPosition[1])
		if state.obstacles[nextKey] {
			state.guardDirection = moveMapping[state.guardDirection]
		}

		state.guardPosition = Position{
			state.guardPosition[0] + positionsMapping[state.guardDirection][0],
			state.guardPosition[1] + positionsMapping[state.guardDirection][1],
		}
	}

	return len(visited)
}

func part2(grid [][]string) int {
	state := initialise(grid)
	startingPosition := state.guardPosition
	startingDirection := state.guardDirection
	visited := make(map[string]bool)

	for {
		key := fmt.Sprintf("%d_%d", state.guardPosition[0], state.guardPosition[1])
		if state.areaLimits[key] {
			visited[key] = true
			break
		}

		visited[key] = true
		nextPosition := Position{
			state.guardPosition[0] + positionsMapping[state.guardDirection][0],
			state.guardPosition[1] + positionsMapping[state.guardDirection][1],
		}

		nextKey := fmt.Sprintf("%d_%d", nextPosition[0], nextPosition[1])
		if state.obstacles[nextKey] {
			state.guardDirection = moveMapping[state.guardDirection]
		}

		state.guardPosition = Position{
			state.guardPosition[0] + positionsMapping[state.guardDirection][0],
			state.guardPosition[1] + positionsMapping[state.guardDirection][1],
		}
	}

	type simulationResult struct {
		cell   string
		isLoop bool
	}
	resultChan := make(chan simulationResult)
	activeSimulations := 0

	for cell := range visited {
		if cell == fmt.Sprintf("%d_%d", startingPosition[0], startingPosition[1]) {
			continue
		}

		activeSimulations++
		go func(cell string) {
			result := simulate(startingPosition, startingDirection, cell, state.areaLimits, state.obstacles)
			resultChan <- simulationResult{cell: cell, isLoop: result}
		}(cell)
	}

	newObstructions := make(map[string]bool)
	for i := 0; i < activeSimulations; i++ {
		result := <-resultChan
		if result.isLoop {
			newObstructions[result.cell] = true
		}
	}

	return len(newObstructions)
}

func readData(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, strings.Split(line, ""))
	}

	return grid
}

func main() {
	grid := readData("data.txt")
	fmt.Printf("Part 1: %d\n", part1(grid))
	fmt.Printf("Part 2: %d\n", part2(grid))
}
