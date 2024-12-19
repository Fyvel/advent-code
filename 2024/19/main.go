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

type Pattern string
type Inputs struct {
	availablePatterns map[Pattern]bool
	desiredPatterns   []Pattern
}

func formatData(rows []string) Inputs {
	inputs := Inputs{
		availablePatterns: map[Pattern]bool{},
		desiredPatterns:   []Pattern{},
	}
	for i, row := range rows {
		if row == "" {
			continue
		}
		if i == 0 {
			availablePatterns := strings.Split(row, ", ")
			for _, pattern := range availablePatterns {
				inputs.availablePatterns[Pattern(pattern)] = true
			}
		} else {
			inputs.desiredPatterns = append(inputs.desiredPatterns, Pattern(row))
		}
	}

	return inputs
}

func (pattern *Pattern) canBeMadeFrom(availablePatterns map[Pattern]bool) bool {
	memo := make(map[int]bool, len(*pattern)+1)
	memo[0] = true

	for rightIndex := 1; rightIndex <= len(*pattern); rightIndex++ {
		// try all substrings
		for leftIndex := 0; leftIndex < rightIndex; leftIndex++ {
			target := Pattern((*pattern)[leftIndex:rightIndex])
			// left is known & target is available
			if memo[leftIndex] && availablePatterns[target] {
				memo[rightIndex] = true
				break
			}
		}
	}

	return memo[len(*pattern)]
}

func part1(inputs Inputs) int {

	possiblePatternsCount := 0

	for _, desiredPattern := range inputs.desiredPatterns {
		if desiredPattern.canBeMadeFrom(inputs.availablePatterns) {
			possiblePatternsCount++
		}
	}

	fmt.Println(possiblePatternsCount)
	return possiblePatternsCount
}

func (pattern *Pattern) calculatePatternCombinations(availablePatterns map[Pattern]bool) int {
	memo := make(map[int]int, len(*pattern)+1)
	memo[0] = 1

	for rightIndex := 1; rightIndex <= len(*pattern); rightIndex++ {
		// try all substrings
		for leftIndex := 0; leftIndex < rightIndex; leftIndex++ {
			target := Pattern((*pattern)[leftIndex:rightIndex])
			// known combination & target is available
			if memo[leftIndex] > 0 && availablePatterns[target] {
				// add combination(s)
				memo[rightIndex] += memo[leftIndex]
			}
		}
	}

	return memo[len(*pattern)]
}

func part2(inputs Inputs) int {

	sumPossiblePatternCombinations := 0

	for _, desiredPattern := range inputs.desiredPatterns {
		sumPossiblePatternCombinations += desiredPattern.calculatePatternCombinations(inputs.availablePatterns)
	}

	fmt.Println(sumPossiblePatternCombinations)
	return sumPossiblePatternCombinations
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
