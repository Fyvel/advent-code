package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func part1(data []string) int {
	multiplierRegex := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)`)
	mulRegex := regexp.MustCompile(`(\d)+`)

	multipliers := multiplierRegex.FindAllString(strings.Join(data, ""), -1)

	sum := 0

	for _, mul := range multipliers {
		nums := mulRegex.FindAllString(mul, -1)
		a, err1 := strconv.Atoi(nums[0])
		b, err2 := strconv.Atoi(nums[1])
		if err1 == nil && err2 == nil {
			sum += a * b
		}

	}
	fmt.Println("sum", sum)
	return sum
}

func part2(data []string) int {
	multiplierRegex := regexp.MustCompile(`(do(n't)?\(\))|(mul\(\d{1,3},\d{1,3}\))`)
	mulRegex := regexp.MustCompile(`\d+`)
	instructionRegex := regexp.MustCompile(`do(n't)?\(\)`)
	sum := 0

	instruction := "do()"

	var multipliers []string
	for _, line := range data {
		matches := multiplierRegex.FindAllString(line, -1)
		multipliers = append(multipliers, matches...)
	}

	for _, mul := range multipliers {
		if instructionRegex.MatchString(mul) {
			instruction = mul
		} else {
			if instruction != "do()" {
				continue
			}
			nums := mulRegex.FindAllString(mul, -1)
			a, err1 := strconv.Atoi(nums[0])
			b, err2 := strconv.Atoi(nums[1])
			if err1 == nil && err2 == nil {
				sum += a * b
			}
		}
	}
	fmt.Println("sum", sum)
	return sum
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData)
	part2(formattedData)
}
