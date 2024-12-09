package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatData(rows []string) ([][]string, error) {
	grid := make([][]string, len(rows))

	for i, row := range rows {
		grid[i] = strings.Split(row, "")
	}

	return grid, nil
}

func visualiseGrid(grid [][]string, antiNodes map[string]bool) {
	var sb strings.Builder

	for i, row := range grid {
		for j, cell := range row {
			key := fmt.Sprintf("%d_%d", j, i)
			if antiNodes[key] {
				sb.WriteString("◆")
			} else {
				sb.WriteString(cell)
			}
		}
		if i < len(grid)-1 {
			sb.WriteString("\n")
		}
	}

	fmt.Println(sb.String())

	if len(grid) > 0 {
		fmt.Println(strings.Repeat("-", len(grid[0])))
	}
}

func isWithinBounds(x, y int, grid [][]string) bool {
	return x >= 0 && y >= 0 && x < len(grid[0]) && y < len(grid)
}

type Location struct{ x, y int }
type AntennasMap map[string][]Location
type AntiNodes map[string]bool

func part1(grid [][]string) int {
	visualiseGrid(grid, make(map[string]bool))

	antennasMap := make(AntennasMap)
	antiNodes := make(AntiNodes)

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			cell := grid[i][j]

			if match, _ := regexp.MatchString(`[a-zA-Z0-9]`, cell); !match {
				continue
			}

			if _, exists := antennasMap[cell]; !exists {
				antennasMap[cell] = []Location{}
			}
			antennasMap[cell] = append(antennasMap[cell], Location{j, i})

			for _, locations := range antennasMap {
				for i := 0; i < len(locations); i++ {
					for j := i + 1; j < len(locations); j++ {
						loc1 := locations[i]
						loc2 := locations[j]

						dx := loc1.x - loc2.x
						dy := loc1.y - loc2.y

						antiNode1 := []int{loc1.x + dx, loc1.y + dy}
						antiNode2 := []int{loc2.x - dx, loc2.y - dy}

						if isWithinBounds(antiNode1[0], antiNode1[1], grid) {
							antiNodes[fmt.Sprintf("%d_%d", antiNode1[0], antiNode1[1])] = true
						}

						if isWithinBounds(antiNode2[0], antiNode2[1], grid) {
							antiNodes[fmt.Sprintf("%d_%d", antiNode2[0], antiNode2[1])] = true
						}
					}
				}
			}
		}
	}

	visualiseGrid(grid, antiNodes)
	fmt.Println(len(antiNodes))
	return len(antiNodes)
}

func part2(grid [][]string) int {
	visualiseGrid(grid, nil)

	antennasMap := make(AntennasMap)
	antiNodes := make(AntiNodes)

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			cell := grid[i][j]

			if match, _ := regexp.MatchString(`[a-zA-Z0-9]`, cell); !match {
				continue
			}

			if _, exists := antennasMap[cell]; !exists {
				antennasMap[cell] = []Location{}
			}
			antennasMap[cell] = append(antennasMap[cell], Location{j, i})

			for _, locations := range antennasMap {
				for i := 0; i < len(locations); i++ {
					for j := i + 1; j < len(locations); j++ {
						loc1 := locations[i]
						loc2 := locations[j]

						dx := loc1.x - loc2.x
						dy := loc1.y - loc2.y

						k := 0
						for {
							localAntiNodes := make(AntiNodes)
							antiNode1 := []int{loc1.x + dx*k, loc1.y + dy*k}
							antiNode2 := []int{loc2.x - dx*k, loc2.y - dy*k}

							if isWithinBounds(antiNode1[0], antiNode1[1], grid) {
								localAntiNodes[fmt.Sprintf("%d_%d", antiNode1[0], antiNode1[1])] = true
							}
							if isWithinBounds(antiNode2[0], antiNode2[1], grid) {
								localAntiNodes[fmt.Sprintf("%d_%d", antiNode2[0], antiNode2[1])] = true
							}
							if len(localAntiNodes) == 0 {
								break
							}
							for node := range localAntiNodes {
								antiNodes[node] = true
							}
							k++
						}
					}
				}
			}
		}
	}

	visualiseGrid(grid, antiNodes)
	fmt.Println(len(antiNodes))
	return len(antiNodes)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData, err := formatData(data)
	part1(formattedData)
	part2(formattedData)
}
