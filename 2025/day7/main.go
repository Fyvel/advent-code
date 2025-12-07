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
			renderGrid(manifoldDiagram, row, col, true)
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
	renderGrid(manifoldDiagram, -1, -1, false)
	fmt.Println("\nPart 1: Processing complete", count)
}

func part2(data [][]string) {
	fmt.Println("Part 2:", data)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData)
	// part2(formattedData)
}

const (
	Reset       = "\033[0m"
	White       = "\033[97m"
	Black       = "\033[30m"
	Yellow      = "\033[33m"
	Pink        = "\033[95m"
	Cyan        = "\033[96m"
	BgCyan      = "\033[46m"
	ClearScreen = "\033[2J"
	MoveCursor  = "\033[H"
)

func renderGrid(grid [][]string, activeRow, activeCol int, clearScreen bool) {
	if clearScreen {
		fmt.Print(ClearScreen + MoveCursor)
	}

	for rowIdx, row := range grid {
		for colIdx, cell := range row {
			isActive := rowIdx == activeRow && colIdx == activeCol

			if isActive {
				fmt.Print(BgCyan + Black)
			}

			switch cell {
			case "S":
				fmt.Print(White + cell + Reset)
			case ".":
				fmt.Print(Black + cell + Reset)
			case "|":
				fmt.Print(Yellow + cell + Reset)
			case "^":
				fmt.Print(Pink + cell + Reset)
			default:
				fmt.Print(cell)
			}

			if isActive {
				fmt.Print(Reset)
			}
		}
		fmt.Println()
	}
	time.Sleep(10 * time.Millisecond)
}
