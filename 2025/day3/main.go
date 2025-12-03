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

func part1(banks []string) {
	sum := 0

	for _, bank := range banks {
		firstDigit := 0
		secondDigit := 0
		firstDigitIdx := -1

		for idx, char := range bank {
			digit, _ := strconv.Atoi(string(char))

			if idx < len(bank)-1 {
				if digit > firstDigit {
					firstDigit = digit
					firstDigitIdx = idx
				} else if digit == firstDigit && idx < firstDigitIdx {
					firstDigitIdx = idx
				}
			}
		}

		rightSide := bank[firstDigitIdx+1:]
		for _, char := range rightSide {
			digit, _ := strconv.Atoi(string(char))
			if digit > secondDigit {
				secondDigit = digit
			}
		}

		highestNumber := firstDigit*10 + secondDigit
		// fmt.Println("Highest number:", highestNumber)
		sum += highestNumber
	}
	fmt.Println("Sum of max values:", sum)
}

func part2(banks []string) {
	sum := 0
	length := 12

	for _, bank := range banks {
		if len(bank) < length {
			continue
		}

		joltage := ""
		startIdx := 0

		for idx := range length {
			remaining := length - idx - 1
			endIdx := len(bank) - remaining

			maxDigit := '0'
			maxDigitIdx := startIdx

			for i := startIdx; i < endIdx; i++ {
				if bank[i] > byte(maxDigit) {
					maxDigit = rune(bank[i])
					maxDigitIdx = i
				}
			}

			startIdx = maxDigitIdx + 1
			joltage += string(maxDigit)
		}

		// fmt.Println("Le joltage:", joltage)
		joltageNum, _ := strconv.Atoi(joltage)
		sum += joltageNum
	}
	fmt.Println("Sum of max values:", sum)
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
