package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

type Vector struct {
	x, y int
}
type CheatSegment struct {
	start Vector
	end   Vector
}
type Maze struct {
	grid      [][]string
	walls     map[Vector]bool
	start     Vector
	end       Vector
	bestPaths []Vector
	cheats    map[CheatSegment]int
}

func formatData(rows []string) Maze {
	maze := Maze{
		grid:      make([][]string, len(rows)),
		walls:     make(map[Vector]bool),
		bestPaths: []Vector{},
		cheats:    make(map[CheatSegment]int),
	}

	for r, row := range rows {
		maze.grid[r] = make([]string, len(row))
		for c, char := range row {
			v := Vector{c, r}
			switch char {
			case '#':
				maze.walls[v] = true
			case 'S':
				maze.start = v
			case 'E':
				maze.end = v
			}
			maze.grid[r][c] = string(char)
		}
	}

	return maze
}

func (maze *Maze) solve() *Maze {
	// North, East, South, West
	directions := []Vector{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	currentTile := maze.start
	maze.bestPaths = append(maze.bestPaths, currentTile)
	lastTile := Vector{-1, -1}

	for maze.grid[currentTile.y][currentTile.x] != "E" {
		tile := Vector{}

		for _, dir := range directions {
			nextTile := Vector{currentTile.x + dir.x, currentTile.y + dir.y}
			if nextTile != lastTile && !maze.walls[nextTile] {
				tile = nextTile
				break
			}
		}
		maze.bestPaths = append(maze.bestPaths, tile)
		lastTile = currentTile
		currentTile = tile
	}

	return maze
}

func (maze *Maze) identifyCheats(maxLength int) *Maze {
	// for each bestPath
	for bestPathIndex, bestPath := range maze.bestPaths {

		// look for next bestPaths within maxLength
		for nextIndex := bestPathIndex; nextIndex < len(maze.bestPaths); nextIndex++ {
			next := maze.bestPaths[nextIndex]

			// calculate distance between current best path & next best path
			distance := abs(next.x-bestPath.x) + abs(next.y-bestPath.y)

			// if distance is within maxLength, it's a cheat/shortcut
			if distance <= maxLength {
				maze.cheats[CheatSegment{bestPath, next}] = nextIndex - bestPathIndex - distance
			}
		}
	}
	return maze
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m *Maze) copy() *Maze {
	newMaze := Maze{
		grid:      make([][]string, len(m.grid)),
		walls:     make(map[Vector]bool),
		start:     m.start,
		end:       m.end,
		bestPaths: make([]Vector, len(m.bestPaths)),
		cheats:    make(map[CheatSegment]int),
	}

	for i := range m.grid {
		newMaze.grid[i] = make([]string, len(m.grid[i]))
		copy(newMaze.grid[i], m.grid[i])
	}
	for k, v := range m.walls {
		newMaze.walls[k] = v
	}
	copy(newMaze.bestPaths, m.bestPaths)
	for k, v := range m.cheats {
		newMaze.cheats[k] = v
	}
	return &newMaze
}

func (m *Maze) renderGrid() {
	// clear screen
	fmt.Print("\033[H\033[2J")

	red := "\033[31m"
	green := "\033[32m"
	orange := "\033[33m"
	reset := "\033[0m"
	fmt.Println()

	bestPathMap := make(map[Vector]bool)
	for _, v := range m.bestPaths {
		bestPathMap[v] = true
	}

	render := make([][]string, len(m.grid))
	for r := range render {
		render[r] = make([]string, len(m.grid[r]))
		for c := range render[r] {
			v := Vector{c, r}
			if m.walls[v] {
				render[r][c] = orange + "#" + reset
			} else if bestPathMap[v] {
				render[r][c] = green + "O" + reset
			} else {
				if v == m.start {
					render[r][c] = "S"
				} else if v == m.end {
					render[r][c] = "E"
				} else {
					render[r][c] = "."
				}
			}
			// check if v is start or end of a cheat
			for cheat := range m.cheats {
				if v == cheat.start || v == cheat.end {
					render[r][c] = red + "#" + reset
				}
			}
		}
		fmt.Print("\r" + strings.Join(render[r], " ") + "\n")
	}
	fmt.Println()
	time.Sleep(250 * time.Millisecond)
}

func part1(input Maze) int {
	maze := input.copy()
	maze.renderGrid()
	maze.solve()
	maze.renderGrid()
	maze.identifyCheats(2)
	maze.renderGrid()

	minPicosecondSaved := 100
	countCheats := 0

	for _, v := range maze.cheats {
		if v >= minPicosecondSaved {
			countCheats++
		}
	}

	fmt.Println("Part 1:", countCheats)
	return countCheats
}

func part2(input Maze) int {
	maze := input.copy()
	// maze.renderGrid()
	maze.solve()
	// maze.renderGrid()
	maze.identifyCheats(20)
	// maze.renderGrid()

	minPicosecondSaved := 100
	countCheats := 0

	for _, v := range maze.cheats {
		if v >= minPicosecondSaved {
			countCheats++
		}
	}

	fmt.Println("Part 2:", countCheats)
	return countCheats
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
