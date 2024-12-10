package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

func bfs(
	grid [][]string,
	startPoints [][]int,
	aggregator func(x, y, val int, result interface{}),
	visited map[string]bool,
	result interface{},
) {

	// init queue
	queue := make([][]int, 0)
	queue = append(queue, startPoints...)

	directions := [][]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]
		x, y := cell[0], cell[1]
		cellValue, err := strconv.Atoi(grid[x][y])
		if err != nil {
			continue
		}
		cellKey := fmt.Sprintf("%d_%d", x, y)
		visited[cellKey] = true

		// process
		aggregator(x, y, cellValue, result)

		for _, dir := range directions {
			newX, newY := x+dir[0], y+dir[1]

			// out of bounds
			if newX < 0 || newY < 0 || newX >= len(grid) || newY >= len(grid[0]) {
				continue
			}

			neighbourValue, err := strconv.Atoi(grid[newX][newY])
			if err != nil {
				continue
			}
			key := fmt.Sprintf("%d_%d", newX, newY)

			// queue up
			if neighbourValue == cellValue+1 {
				queue = append(queue, []int{newX, newY})
				visited[key] = true
			}
		}
	}
}

func part1(grid [][]string) int {

	startPoints := make([][]int, 0)
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == fmt.Sprint(0) {
				startPoints = append(startPoints, []int{i, j})
			}
		}
	}

	trailHeadScores := make(map[string]int)
	for _, startPoint := range startPoints {

		cellVisited := make(map[string]bool)
		summitVisited := make(map[string]bool)

		scoresAggregator := func(x, y, value int, result interface{}) {

			key := fmt.Sprintf("%d_%d", x, y)

			if value == 9 && !summitVisited[key] {
				trailHeadKey := fmt.Sprintf("%d_%d", startPoint[0], startPoint[1])
				summitVisited[key] = true
				result.(map[string]int)[trailHeadKey]++
			}
		}

		bfs(grid, [][]int{startPoint}, scoresAggregator, cellVisited, trailHeadScores)
	}

	sum := 0
	for _, score := range trailHeadScores {
		sum += score
	}
	fmt.Println(sum)
	return sum
}

func part2(grid [][]string) int {

	startPoints := make([][]int, 0)
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == fmt.Sprint(0) {
				startPoints = append(startPoints, []int{i, j})
			}
		}
	}

	trailHeadScores := make(map[string]int)
	for _, startPoint := range startPoints {

		cellVisited := make(map[string]bool)
		summitVisited := make(map[string]bool)

		scoresAggregator := func(x, y, value int, result interface{}) {

			key := fmt.Sprintf("%d_%d", x, y)

			if value == 9 {
				trailHeadKey := fmt.Sprintf("%d_%d", startPoint[0], startPoint[1])
				summitVisited[key] = true
				result.(map[string]int)[trailHeadKey]++
			}
		}

		bfs(grid, [][]int{startPoint}, scoresAggregator, cellVisited, trailHeadScores)
	}

	sum := 0
	for _, score := range trailHeadScores {
		sum += score
	}
	fmt.Println(sum)
	return sum
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData, err := formatData(data)
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}
	part1(formattedData)
	part2(formattedData)
}
