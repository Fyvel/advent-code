package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func readDataAlt() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatDataAlt(rows []string) ([][]string, error) {
	grid := make([][]string, len(rows))

	for i, row := range rows {
		grid[i] = strings.Split(row, "")
	}

	return grid, nil
}

func bfsAlt(
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

func part1Routine(grid [][]string) int {

	startPoints := make([][]int, 0)
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == fmt.Sprint(0) {
				startPoints = append(startPoints, []int{i, j})
			}
		}
	}

	var syncLock sync.Mutex
	var waitGroup sync.WaitGroup

	trailHeadScores := make(map[string]int)
	for _, startPoint := range startPoints {

		waitGroup.Add(1)

		go func(startPointCell []int) {
			defer waitGroup.Done()

			cellVisited := make(map[string]bool)
			summitVisited := make(map[string]bool)

			scoresAggregator := func(x, y, value int, result interface{}) {
				key := fmt.Sprintf("%d_%d", x, y)

				if value == 9 && !summitVisited[key] {
					trailHeadKey := fmt.Sprintf("%d_%d", startPointCell[0], startPointCell[1])
					summitVisited[key] = true
					syncLock.Lock()
					result.(map[string]int)[trailHeadKey]++
					syncLock.Unlock()
				}
			}

			bfsAlt(grid, [][]int{startPointCell}, scoresAggregator, cellVisited, trailHeadScores)

		}(startPoint)
	}

	waitGroup.Wait()

	sum := 0
	for _, score := range trailHeadScores {
		sum += score
	}
	fmt.Println(sum)
	return sum
}

func part2Routine(grid [][]string) int {
	startPoints := make([][]int, 0)
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == fmt.Sprint(0) {
				startPoints = append(startPoints, []int{i, j})
			}
		}
	}

	var syncLock sync.Mutex
	trailHeadScores := make(map[string]int)
	var waitGroup sync.WaitGroup

	for _, startPoint := range startPoints {
		waitGroup.Add(1)
		go func(sp []int) {
			defer waitGroup.Done()

			cellVisited := make(map[string]bool)

			scoresAggregator := func(x, y, value int, result interface{}) {

				if value == 9 {
					trailHeadKey := fmt.Sprintf("%d_%d", sp[0], sp[1])
					syncLock.Lock()
					result.(map[string]int)[trailHeadKey]++
					syncLock.Unlock()
				}
			}

			bfsAlt(grid, [][]int{sp}, scoresAggregator, cellVisited, trailHeadScores)

		}(startPoint)
	}

	waitGroup.Wait()

	sum := 0
	for _, score := range trailHeadScores {
		sum += score
	}
	fmt.Println(sum)
	return sum
}

func main() {
	data, err := readDataAlt()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData, err := formatDataAlt(data)
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}
	part1Routine(formattedData)
	part2Routine(formattedData)
}
