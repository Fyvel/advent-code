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

func formatData(rows []string) [][]string {
	formatted := make([][]string, len(rows))
	for i, row := range rows {
		formatted[i] = strings.Split(row, "")
	}
	return formatted
}

func part1(data [][]string) {
	manifoldDiagram := data
	count := 0

	beamRange := [2]int{0, len(manifoldDiagram) - 1}

	for row := range manifoldDiagram {
		for col := beamRange[0]; col <= beamRange[1]; col++ {
			// renderGrid(manifoldDiagram, row, col, nil)
			isBeam := false

			if manifoldDiagram[row][col] == "S" {
				beamRange = [2]int{col, col}
				break
			}

			if row < 0 || row >= len(manifoldDiagram) || col < 0 || col >= len(manifoldDiagram[0]) {
				continue
			}

			if manifoldDiagram[row][col] == "^" && manifoldDiagram[row-1][col] == "|" {
				isBeam = true
				manifoldDiagram[row][col-1] = "|"
				manifoldDiagram[row][col+1] = "|"
				if beamRange[0] > col-1 {
					beamRange[0] = col - 1
				}
				if beamRange[1] < col+1 {
					beamRange[1] = col + 1
				}
				count++
				continue
			}

			if isBeam || manifoldDiagram[row][col] == "S" || row == 0 {
				continue
			}

			if manifoldDiagram[row-1][col] == "S" {
				manifoldDiagram[row][col] = "|"
				beamRange = [2]int{col, col}
				continue
			}
			if manifoldDiagram[row-1][col] == "|" && manifoldDiagram[row][col] != "^" {
				manifoldDiagram[row][col] = "|"
				continue
			}
		}
	}

	fmt.Print(ClearScreen + MoveCursor)
	renderGrid(manifoldDiagram, -1, -1, nil)
	fmt.Println("\nPart 1: Processing complete", count)
}

func part2(data [][]string) {
	manifoldDiagram := data

	var startRow, startCol int
	for c := range manifoldDiagram[0] {
		if manifoldDiagram[0][c] == "S" {
			startRow = 0
			startCol = c
			break
		}
	}

	totalPaths := dfsCountPath(manifoldDiagram, startRow+1, startCol, make(map[string]int), make(map[string]bool), make(map[string]bool))

	fmt.Print(ClearScreen + MoveCursor)
	renderGrid(manifoldDiagram, -1, -1, make(map[string]bool))
	fmt.Println("\nPart 2:", totalPaths)
}

func dfsCountPath(grid [][]string, row, col int, memo map[string]int, visited map[string]bool, activePath map[string]bool) int {
	key := fmt.Sprintf("%d_%d", row, col)

	activePath[key] = true
	renderGrid(grid, row, col, activePath)

	// base
	if row >= len(grid) {
		activePath[key] = false
		return 1
	}
	if col < 0 || col >= len(grid[0]) {
		activePath[key] = false
		return 0
	}
	if val, exists := memo[key]; exists {
		activePath[key] = false
		return val
	}

	activeCell := grid[row][col]
	if !visited[key] {
		visited[key] = true
		if activeCell == "." {
			grid[row][col] = "⏐"
		}
	}

	// recursive
	cell := activeCell
	paths := 0
	switch cell {
	case "^":
		leftPaths := dfsCountPath(grid, row+1, col-1, memo, visited, activePath)
		rightPaths := dfsCountPath(grid, row+1, col+1, memo, visited, activePath)
		paths = leftPaths + rightPaths
	default:
		paths = dfsCountPath(grid, row+1, col, memo, visited, activePath)
	}

	memo[key] = paths
	activePath[key] = false
	return paths
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	part1(formatData(data))

	part2(formatData(data))
}

const (
	Reset        = "\033[0m"
	White        = "\033[97m"
	Black        = "\033[30m"
	Orange       = "\033[38;5;208m"
	Cyan         = "\033[96m"
	BgOrange     = "\033[48;5;208m"
	ClearScreen  = "\033[2J"
	MoveCursor   = "\033[H"
	AltScreenOn  = "\033[?1049h"
	AltScreenOff = "\033[?1049l"
	HideCursor   = "\033[?25l"
	ShowCursor   = "\033[?25h"
)

func renderGrid(grid [][]string, activeRow, activeCol int, activePath map[string]bool) {
	var buf strings.Builder
	buf.WriteString(MoveCursor)

	for rowIdx, row := range grid {
		for colIdx, cell := range row {
			isActive := rowIdx == activeRow && colIdx == activeCol
			key := fmt.Sprintf("%d_%d", rowIdx, colIdx)
			isInActivePath := activePath != nil && activePath[key]

			if isActive {
				buf.WriteString(BgOrange + Black)
			} else if isInActivePath {
				buf.WriteString(Orange)
			}

			switch cell {
			case "S", "^":
				buf.WriteString(White + cell + Reset)
			case ".":
				buf.WriteString(Black + cell + Reset)
			case "⏐": // this is the weird one
				if isInActivePath || isActive {
					buf.WriteString(Orange + cell + Reset)
				} else {
					buf.WriteString(Cyan + cell + Reset)
				}
			case "|":
				buf.WriteString(Cyan + cell + Reset)
			default:
				buf.WriteString(cell)
			}

			if isActive || isInActivePath {
				buf.WriteString(Reset)
			}
		}
		buf.WriteString("\n")
	}

	fmt.Print(buf.String())
	time.Sleep(5 * time.Millisecond)
}
