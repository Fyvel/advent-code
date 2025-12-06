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

func formatData(rows []string) []string {
	return rows
}

func part1(data []string) {
	operators := strings.Fields(data[len(data)-1:][0])
	numbersRows := data[:len(data)-1]

	operatorsMap := make(map[int]string)
	exams := make(map[int]int)

	for i, op := range operators {
		operatorsMap[i] = op
		if op == "*" {
			exams[i] = 1
		} else {
			exams[i] = 0
		}
	}

	for _, numbersRow := range numbersRows {
		fields := strings.Fields(numbersRow)

		for c, field := range fields {
			number, _ := strconv.Atoi(field)
			switch operatorsMap[c] {
			case "*":
				exams[c] *= number
			case "+":
				exams[c] += number
			}
		}
	}

	sum := 0
	for _, val := range exams {
		sum += val
	}
	fmt.Println("Solution:", sum)
}

func part2(data []string) {
	// fmt.Println("Part 2:", data)
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
