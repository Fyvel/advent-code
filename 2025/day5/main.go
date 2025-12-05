package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func readData() ([2][]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return [2][]string{}, err
	}
	result := strings.Split(string(data), "\n\n")
	part1 := strings.Split(result[0], "\n")
	part2 := strings.Split(result[1], "\n")
	return [2][]string{part1, part2}, nil
}

func part1(intervals []string, ingredients []string) {
	freshCount := 0

	mergedIntervals := [][2]int{}
	for _, interval := range intervals {
		parts := strings.Split(interval, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])

		newInterval := [2]int{start, end}
		merged := false

		for i, existing := range mergedIntervals {
			// check for overlap
			if !(newInterval[1] < existing[0] || newInterval[0] > existing[1]) {
				mergedIntervals[i][0] = min(mergedIntervals[i][0], newInterval[0])
				mergedIntervals[i][1] = max(mergedIntervals[i][1], newInterval[1])
				merged = true
				break
			}
		}

		if !merged {
			mergedIntervals = append(mergedIntervals, newInterval)
		}
	}

	for _, ingredient := range ingredients {
		ingredientId, _ := strconv.Atoi(ingredient)

		for _, interval := range mergedIntervals {
			// check if it's in range
			if ingredientId >= interval[0] && ingredientId <= interval[1] {
				freshCount++
				break
			}
		}
	}

	fmt.Println("Solution:", freshCount)
}

// func part2(intervals []string, ingredients []string) {}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	part1(data[0], data[1])
	// part2(data[0], data[1])
}
