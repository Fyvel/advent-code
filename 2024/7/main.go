package main

import (
	"fmt"
	"math"
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

func concatInts(a, b int) int {
	multiplier := int(math.Pow10(int(math.Log10(float64(b))) + 1))
	return a*multiplier + b
}

func parseNum(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func formatData(rows []string) [][][]int {
	equations := make([][][]int, len(rows))

	for i, row := range rows {
		parts := strings.Split(row, ": ")
		numbers := strings.Split(parts[1], " ")

		target := []int{parseNum(parts[0])}

		rest := make([]int, len(numbers))
		for j, num := range numbers {
			rest[j] = parseNum(num)
		}

		equations[i] = [][]int{target, rest}
	}

	return equations
}

func calculate(eq [][]int, operators []func(int, int) int, resultChan ...chan int) int {
	target := eq[0][0]
	numbers := eq[1]

	evaluations := make(map[int]bool)
	evaluations[numbers[0]] = true
	rest := numbers[1:]

	for i := 0; i < len(rest); i++ {
		localEvaluations := make(map[int]bool)

		for left := range evaluations {
			for _, operator := range operators {
				localEvaluations[operator(left, rest[i])] = true
			}
		}

		evaluations = localEvaluations
	}

	if evaluations[target] {
		if len(resultChan) > 0 {
			resultChan[0] <- target
		}
		return target
	}
	if len(resultChan) > 0 {
		resultChan[0] <- 0
	}
	return 0
}

func evaluate(equations [][][]int, operators []func(int, int) int) int {
	sum := 0
	for _, equation := range equations {
		sum += calculate(equation, operators)
	}
	return sum
}

func evaluateParallel(equations [][][]int, operators []func(int, int) int) int {
	resultChan := make(chan int, len(equations))

	for _, equation := range equations {
		go calculate(equation, operators, resultChan)
	}

	sum := 0
	for i := 0; i < len(equations); i++ {
		sum += <-resultChan
	}

	return sum
}

func part1(equations [][][]int) int {
	operators := []func(int, int) int{
		func(a, b int) int { return a + b },
		func(a, b int) int { return a * b },
	}

	sum := evaluateParallel(equations, operators)

	fmt.Printf("Sum: %d\n", sum)
	return sum
}

func part2(equations [][][]int) int {
	operators := []func(int, int) int{
		func(a, b int) int { return a + b },
		func(a, b int) int { return a * b },
		func(a, b int) int { return concatInts(a, b) },
	}

	sum := evaluateParallel(equations, operators)

	fmt.Printf("Sum: %d\n", sum)
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
