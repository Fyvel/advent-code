package main

import (
	"aoc2025/utils"
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
	grid := make([][]string, len(rows))

	for r, c := range rows {
		grid[r] = strings.Split(c, "")
	}

	return grid
}

func part1(grid [][]string) {
	totalRolls := 0

	for r := range grid {
		for c := range grid[r] {

			coords := [2]int{r, c}
			cell := grid[r][c]

			if cell != "@" {
				continue
			}

			directions := [][]int{{0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}}
			rollsCount := 0

			for _, direction := range directions {
				row := coords[0] + direction[0]
				col := coords[1] + direction[1]

				if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[0]) {
					continue
				}

				if grid[row][col] == "@" {
					rollsCount++
				}
			}
			if rollsCount < 4 {
				totalRolls++
			}
		}
	}

	fmt.Println("Part 1:", totalRolls)

}

func part2(grid [][]string, withVisual bool) {
	totalRolls := 0
	lastCount := -1

	if withVisual {
		utils.EnterVisualMode()
		fmt.Print(utils.ClearScreen + utils.MoveCursor)
	}

	for lastCount != totalRolls {
		lastCount = totalRolls

		if withVisual {
			renderGrid(grid)
		}
		for r := range grid {
			for c := range grid[r] {

				coords := [2]int{r, c}
				cell := grid[r][c]

				if cell != "@" {
					continue
				}

				directions := [][]int{{0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}}
				rollsCount := 0

				for _, direction := range directions {
					row := coords[0] + direction[0]
					col := coords[1] + direction[1]

					if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[0]) {
						continue
					}

					if grid[row][col] == "@" {
						rollsCount++
					}
				}
				if rollsCount < 4 {
					totalRolls++
					grid[r][c] = "x"
				}
			}
		}

		if withVisual {
			renderGrid(grid)
			// replace `x` with `.` for smoother visual
			for r := range grid {
				for c := range grid[r] {
					if grid[r][c] == "x" {
						grid[r][c] = "."
					}
				}
			}
		}
	}
	if withVisual {
		utils.ExitVisualMode()
	}

	fmt.Println("Part 2:", totalRolls)
}

func renderGrid(grid [][]string) {
	cellRenderer := func(ctx utils.CellRenderContext) string {
		switch ctx.Cell {
		case ".":
			return utils.Blue + ctx.Cell + utils.Reset
		case "x":
			return utils.Yellow + ctx.Cell + utils.Reset
		case "@":
			return utils.HotPink + ctx.Cell + utils.Reset
		default:
			return utils.BgBlack + utils.White + ctx.Cell + utils.Reset
		}
	}

	utils.RenderGrid(grid, -1, -1, nil, cellRenderer)
	time.Sleep(100 * time.Millisecond)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	withVisual := os.Getenv("AOC_VISUAL") == "1"

	formattedData := formatData(data)
	part1(formattedData)
	part2(formattedData, withVisual)
}
