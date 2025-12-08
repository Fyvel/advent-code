package main

import (
	"aoc2025/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
			isBeam := false
			// utils.RenderGrid(manifoldDiagram, row, col, nil, cellRenderer)

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

	fmt.Print(utils.ClearScreen + utils.MoveCursor)
	utils.RenderGrid(manifoldDiagram, -1, -1, nil, nil)
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

	frameCounter := 0
	render := func(grid [][]string, row, col int, activePath map[string]bool) {
		if frameCounter%2 == 0 {
			utils.RenderGrid(grid, row, col, activePath, cellRenderer)
		}
		frameCounter++
	}

	totalPaths := dfsCountPath(manifoldDiagram, startRow+1, startCol, make(map[string]int), make(map[string]bool), make(map[string]bool), render)

	fmt.Print(utils.ClearScreen + utils.MoveCursor)
	utils.RenderGrid(manifoldDiagram, -1, -1, nil, cellRenderer)
	fmt.Println("\nPart 2:", totalPaths)
}

type Render func(grid [][]string, row, col int, activePath map[string]bool)

func dfsCountPath(grid [][]string, row, col int, memo map[string]int, visited map[string]bool, activePath map[string]bool, render Render) int {
	key := fmt.Sprintf("%d_%d", row, col)

	// base
	if row >= len(grid) {
		return 1
	}
	if col < 0 || col >= len(grid[0]) {
		return 0
	}
	if val, exists := memo[key]; exists {
		return val
	}

	activePath[key] = true
	if render != nil {
		render(grid, row, col, activePath)
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
		leftPaths := dfsCountPath(grid, row+1, col-1, memo, visited, activePath, render)
		rightPaths := dfsCountPath(grid, row+1, col+1, memo, visited, activePath, render)
		paths = leftPaths + rightPaths
	default:
		paths = dfsCountPath(grid, row+1, col, memo, visited, activePath, render)
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

func cellRenderer(ctx utils.CellRenderContext) string {
	var buf strings.Builder

	if ctx.IsActive {
		buf.WriteString(utils.BgOrange + utils.Black)
	} else if ctx.IsInActivePath {
		buf.WriteString(utils.Orange)
	}

	switch ctx.Cell {
	case "S", "^":
		buf.WriteString(utils.White + ctx.Cell + utils.Reset)
	case ".":
		buf.WriteString(utils.Black + ctx.Cell + utils.Reset)
	case "⏐": // this is the weird one
		if ctx.IsInActivePath || ctx.IsActive {
			buf.WriteString(utils.BgOrange + utils.Orange + ctx.Cell + utils.Reset)
		} else {
			buf.WriteString(utils.BgCyan + utils.Cyan + ctx.Cell + utils.Reset)
		}
	case "|":
		buf.WriteString(utils.BgCyan + utils.Cyan + ctx.Cell + utils.Reset)
	default:
		buf.WriteString(ctx.Cell)
	}

	if ctx.IsActive || ctx.IsInActivePath {
		buf.WriteString(utils.Reset)
	}

	return buf.String()
}
