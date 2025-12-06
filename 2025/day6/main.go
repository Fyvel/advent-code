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
	operators := strings.Fields(data[len(data)-1])
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
	operators := strings.Fields(data[len(data)-1])
	rows := data[:len(data)-1]

	exams := make([]int, len(operators))
	for i, op := range operators {
		if op == "*" {
			exams[i] = 1
		}
	}

	examDigits := make([]int, len(rows[0]))
	multipliers := make([]int, len(rows[0]))
	for i := range multipliers {
		multipliers[i] = 1
	}

	for rowIdx := len(rows) - 1; rowIdx >= 0; rowIdx-- {
		row := rows[rowIdx]
		for col := 0; col < len(row); col++ {
			if row[col] >= '0' && row[col] <= '9' {
				digit := int(row[col] - '0')
				examDigits[col] += digit * multipliers[col]
				multipliers[col] *= 10
			}
		}
	}

	opIdx := len(operators) - 1
	for col := len(examDigits) - 1; col >= 0; col-- {
		// check for next exam
		if multipliers[col] == 1 {
			opIdx--
			continue
		}

		switch operators[opIdx] {
		case "*":
			exams[opIdx] *= examDigits[col]
		case "+":
			exams[opIdx] += examDigits[col]
		}
	}

	sum := 0
	for _, val := range exams {
		sum += val
	}
	fmt.Println("Solution 2:", sum)
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
