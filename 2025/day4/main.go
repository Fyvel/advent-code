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

	fmt.Println("Total rolls in reach:", totalRolls)

}

func part2(grid [][]string) {
	totalRolls := 0
	lastCount := -1

	for lastCount != totalRolls {
		lastCount = totalRolls

		renderGrid(grid)
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

		renderGrid(grid)
		// replace `x`` with `.`` for smoother visual
		for r := range grid {
			for c := range grid[r] {
				if grid[r][c] == "x" {
					grid[r][c] = "."
				}
			}
		}
	}
	fmt.Println("Total rolls in reach:", totalRolls)

}

func renderGrid(grid [][]string) {
	blue := "\033[34m"
	yellow := "\033[33m"
	hotpink := "\033[38;5;205m"
	reset := "\033[0m"

	for r := range grid {
		for c := range grid[r] {
			cell := grid[r][c]
			switch cell {
			case ".":
				fmt.Print(blue + cell + reset)
			case "x":
				fmt.Print(yellow + cell + reset)
			case "@":
				fmt.Print(hotpink + cell + reset)
			default:
				fmt.Print(cell)
			}
		}
		fmt.Println()
	}
	fmt.Println()
	time.Sleep(200 * time.Millisecond)
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
