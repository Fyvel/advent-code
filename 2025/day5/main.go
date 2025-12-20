package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

	intervalsList := make([]Interval, 0, len(intervals))
	for _, interval := range intervals {
		parts := strings.Split(interval, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		intervalsList = append(intervalsList, Interval{Start: start, End: end})
	}

	mergedIntervals := mergeIntervals(intervalsList)

	for _, ingredient := range ingredients {
		ingredientId, _ := strconv.Atoi(ingredient)

		for _, interval := range mergedIntervals {
			// check if it's in range
			if ingredientId >= interval.Start && ingredientId <= interval.End {
				freshCount++
				break
			}
		}
	}

	fmt.Println("Part 1:", freshCount)
}

type Interval struct{ Start, End int }

func mergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}

	// order by start
	sort.Slice(intervals, func(a, b int) bool { return intervals[a].Start < intervals[b].Start })

	merged := make([]Interval, 0, len(intervals))
	curr := intervals[0]

	for _, next := range intervals[1:] {
		// extend on overlap or adjacent
		if next.Start <= curr.End {
			if next.End > curr.End {
				curr.End = next.End
			}
			continue
		}
		// concat on non-overlapping
		merged = append(merged, curr)
		curr = next
	}

	// add last
	merged = append(merged, curr)
	return merged
}

func part2(intervals []string) {
	idCount := 0

	intervalsList := make([]Interval, 0, len(intervals))
	for _, interval := range intervals {
		parts := strings.Split(interval, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		intervalsList = append(intervalsList, Interval{Start: start, End: end})
	}

	mergedIntervals := mergeIntervals(intervalsList)

	for _, interval := range mergedIntervals {
		idCount += (interval.End - interval.Start + 1)
	}

	fmt.Println("Part 2:", idCount)

}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	part1(data[0], data[1])
	part2(data[0])
}
