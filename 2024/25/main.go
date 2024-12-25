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

type Lock struct {
	heights [5]int
}
type Key struct {
	heights [5]int
}
type Inputs struct {
	locks []Lock
	keys  []Key
}

func formatData(lines []string) Inputs {
	var result Inputs

	// 7 lines per lock/key + 1 empty line
	for i := 0; i < len(lines); i += 8 {
		if i+7 > len(lines) {
			break
		}

		current := [5]int{}
		isLock := false

		for row := 0; row < 7; row++ {
			for c, char := range lines[i+row] {
				if char == '#' {
					current[c]++
				}
			}
			if row == 0 && current == [5]int{1, 1, 1, 1, 1} {
				isLock = true
			}
		}

		// deduct the first row (lock/key identifier)
		for i := range current {
			current[i]--
		}

		if isLock {
			result.locks = append(result.locks, Lock{heights: current})
		} else {
			result.keys = append(result.keys, Key{heights: current})
		}
	}
	return result
}

func part1(inputs Inputs) {
	count := 0

	for _, lock := range inputs.locks {
		for _, key := range inputs.keys {

			isMatching := true
			for i := 0; i < 5; i++ {
				if lock.heights[i]+key.heights[i] > 5 {
					isMatching = false
					break
				}
			}

			if isMatching {
				count++
			}
		}
	}

	fmt.Println("Part 1:", count)
}

// func part2() {}

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
