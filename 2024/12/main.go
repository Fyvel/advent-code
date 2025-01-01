package main

import (
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
	grid := make([][]string, len(rows))

	for i, row := range rows {
		grid[i] = strings.Split(row, "")
	}

	return grid
}

type Area struct {
	region    string
	size      int
	perimeter int
	corners   int
}

func exploreRegion(grid [][]string, region string, start [2]int, visited map[string]bool) Area {
	directions := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	queue := [][2]int{start}

	visited[fmt.Sprintf("%d_%d", start[0], start[1])] = true

	area := Area{
		region:    region,
		size:      1,
		perimeter: 0,
		corners:   0,
	}

	areaCorners := make(map[[2]int]int)

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]

		for _, direction := range directions {
			row := cell[0] + direction[0]
			col := cell[1] + direction[1]
			areaCorners[cell] = countCellCorners(grid, cell)

			// within bounds
			if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[0]) || grid[row][col] != region {
				area.perimeter++
				continue
			}

			key := fmt.Sprintf("%d_%d", row, col)

			if !visited[key] && grid[row][col] == region {
				queue = append(queue, [2]int{row, col})
				visited[key] = true
				area.size++
			}
		}
	}

	for _, corners := range areaCorners {
		area.corners += corners
	}

	return area
}

func part1(grid [][]string) int {
	visited := make(map[string]bool)
	var areas []Area

	// traverse grid
	for r := range grid {
		for c := range grid[r] {
			key := fmt.Sprintf("%d_%d", r, c)
			if visited[key] {
				continue
			}

			region := grid[r][c]
			area := exploreRegion(grid, region, [2]int{r, c}, visited)

			areas = append(areas, area)
		}
	}

	sum := 0
	for _, area := range areas {
		sum += area.size * area.perimeter
	}

	fmt.Println(sum)
	return sum
}

func countCellCorners(grid [][]string, tile [2]int) int {
	corners := 0
	region := grid[tile[0]][tile[1]]

	// W WN N NE E ES S SW
	directions := [][]int{{0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}}

	neighbors := make([]string, 8)
	for i, direction := range directions {
		row := tile[0] + direction[0]
		col := tile[1] + direction[1]

		if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[0]) {
			neighbors[i] = ""
			continue
		}

		neighbors[i] = grid[row][col]
	}

	angleAndAdjacent := [4]map[string]string{}

	angleAndAdjacent[0] = map[string]string{
		"adjacentLeft":  neighbors[0],
		"angle":         neighbors[1],
		"adjacentRight": neighbors[2],
	}
	angleAndAdjacent[1] = map[string]string{
		"adjacentLeft":  neighbors[2],
		"angle":         neighbors[3],
		"adjacentRight": neighbors[4],
	}
	angleAndAdjacent[2] = map[string]string{
		"adjacentLeft":  neighbors[4],
		"angle":         neighbors[5],
		"adjacentRight": neighbors[6],
	}
	angleAndAdjacent[3] = map[string]string{
		"adjacentLeft":  neighbors[6],
		"angle":         neighbors[7],
		"adjacentRight": neighbors[0],
	}

	for _, angle := range angleAndAdjacent {
		// fmt.Println(tile, idx, angle)

		// outer corners (corner and adjacent are different regions)
		if angle["angle"] != region && angle["adjacentLeft"] != region && angle["adjacentRight"] != region {
			corners++
		}

		// inner corners (corner is same region and adjacent are different region)
		if (angle["angle"] == region && angle["adjacentLeft"] != region && angle["adjacentRight"] != region) || (angle["angle"] != region && angle["adjacentLeft"] == region && angle["adjacentRight"] == region) {
			corners++
		}
	}

	return corners
}

func part2(grid [][]string) int {
	visited := make(map[string]bool)
	var areas []Area

	for r := range grid {
		for c := range grid[r] {
			key := fmt.Sprintf("%d_%d", r, c)
			if visited[key] {
				continue
			}

			region := grid[r][c]
			area := exploreRegion(grid, region, [2]int{r, c}, visited)
			areas = append(areas, area)
		}
	}

	sumSides := 0
	for _, area := range areas {
		sumSides += area.size * area.corners
	}

	fmt.Println(sumSides)
	return sumSides
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData)
	part2(formattedData)
}
