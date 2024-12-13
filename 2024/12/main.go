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
	sides     int
}

func exploreRegion(grid [][]string, region string, start []int, visited map[string]bool) Area {
	directions := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	queue := [][]int{start}

	visited[fmt.Sprintf("%d,%d", start[0], start[1])] = true

	area := Area{
		region:    region,
		size:      1,
		perimeter: 0,
	}

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]

		for _, direction := range directions {
			row := cell[0] + direction[0]
			col := cell[1] + direction[1]

			// within bounds
			if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[start[0]]) {
				area.perimeter++
				if direction[0] != lastDirection[0] || direction[1] != lastDirection[1] {
					area.sides++
					lastDirection = direction
				}
				continue
			}
			if grid[row][col] != region {
				area.perimeter++
				if direction[0] != lastDirection[0] || direction[1] != lastDirection[1] {
					area.sides++
					lastDirection = direction
				}
				continue
			}

			key := fmt.Sprintf("%d,%d", row, col)
			area.size++

			if !visited[key] && grid[row][col] == region {
				queue = append(queue, []int{row, col})
				visited[key] = true
				area.size++
				continue
			}
		}
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
			area := exploreRegion(grid, region, []int{r, c}, visited)

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

func getRegion(grid [][]string, cell []int) string {
	if cell[0] < 0 || cell[0] >= len(grid) || cell[1] < 0 || cell[1] >= len(grid[cell[0]]) {
		return ""
	}
	return grid[cell[1]][cell[0]]
}

func countCorners(grid [][]string, cell []int, directions [][]int, dirIndex int, region string) int {
	corners := 0

	// TODO: fix this!
	for d := dirIndex; d < 4; d++ {
		// get current & next directions
		dx1, dy1 := directions[d][0], directions[d][1]
		dx2, dy2 := directions[(d+1)%4][0], directions[(d+1)%4][1]

		left := getRegion(grid, []int{cell[1] + dy1, cell[0] + dx1})
		right := getRegion(grid, []int{cell[1] + dy2, cell[0] + dx2})
		mid := getRegion(grid, []int{cell[1] + dy1 + dy2, cell[0] + dx1 + dx2})

		if (left != region && right != region) || (left == region && right == region && mid != region) {
			corners++
		}
	}

	return corners
}

func exploreRegion2(grid [][]string, region string, start []int, visited map[string]bool) Area {
	directions := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	queue := [][]int{start}

	visited[fmt.Sprintf("%d_%d", start[0], start[1])] = true

	area := Area{
		region:    region,
		size:      1,
		perimeter: 0,
		sides:     0,
	}

	lastEdgeDirIndex := -1

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]
		edges := map[string]map[int]bool{}

		for dirIndex, direction := range directions {
			row := cell[0] + direction[0]
			col := cell[1] + direction[1]
			nextCellKey := fmt.Sprintf("%d_%d", row, col)

			// within bounds
			if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[start[0]]) {
				area.perimeter++
				edges[nextCellKey] = map[int]bool{dirIndex: true}
				if lastEdgeDirIndex != dirIndex {
					area.sides += countCorners(grid, cell, directions, dirIndex, region)
					lastEdgeDirIndex = dirIndex
				}
				continue
			}
			if grid[row][col] != region {
				area.perimeter++
				edges[nextCellKey] = map[int]bool{dirIndex: true, direction[1]: true}
				if lastEdgeDirIndex != dirIndex {
					area.sides += countCorners(grid, cell, directions, dirIndex, region)
					lastEdgeDirIndex = dirIndex
				}
				continue
			}

			key := fmt.Sprintf("%d_%d", row, col)

			if !visited[key] && grid[row][col] == region {
				queue = append(queue, []int{row, col})
				visited[key] = true
				area.size++
				continue
			}
		}
	}

	return area
}

func part2(grid [][]string) int {
	visited := make(map[string]bool)
	var areas []Area

	// traverse grid
	for r := range grid {
		for c := range grid[r] {
			key := fmt.Sprintf("%d_%d", r, c)
			cell := []int{r, c}
			if visited[key] {
				continue
			}

			region := grid[r][c]
			area := exploreRegion2(grid, region, cell, visited)
			areas = append(areas, area)
		}
	}

	sumSides := 0
	sum := 0
	for _, area := range areas {
		sumSides += area.size * area.sides
		sum += area.size * area.perimeter
	}

	fmt.Println(areas)
	fmt.Println(sum)
	fmt.Println(sumSides)
	return sumSides
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
