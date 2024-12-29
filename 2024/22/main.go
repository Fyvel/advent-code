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

func formatData(rows []string) []int {
	data := make([]int, len(rows))
	for i, row := range rows {
		num, _ := strconv.Atoi(row)
		data[i] = num
	}
	return data
}

func simulateSecret(initial int, iterations int) int {
	secret := initial
	for i := 0; i < iterations; i++ {
		secret = generateSecret(secret)
	}
	return secret
}

func part1(data []int) {
	totalSum := 0
	for _, initial := range data {
		finalSecret := simulateSecret(initial, 2000)
		totalSum += finalSecret
	}
	fmt.Println("Total Sum:", totalSum)
}

func generateSecret(secret int) int {
	// multiply by 64 and mix & prune
	secret ^= (secret * 64)
	secret %= 16777216

	// divide by 32 and mix & prune
	secret ^= (secret / 32)
	secret %= 16777216

	// multiply by 2048 and mix & prune
	secret ^= (secret * 2048)
	secret %= 16777216
	return secret
}

const windowSize = 4

type Window struct {
	digits [windowSize]int
}

func getSoldBananas(secret int) []int {
	results := make([]int, 0)
	currentSecret := secret

	for i := 0; i < 2000; i++ {
		newSecret := generateSecret(currentSecret)
		results = append(results, newSecret%10)
		currentSecret = newSecret
	}
	return results
}

func calculateDifferences(numbers []int) []int {
	length := len(numbers) - 1
	diff := make([]int, length-1)
	for i := 0; i < length-1; i++ {
		diff[i] = numbers[i+1] - numbers[i]
	}
	return diff
}

func getBuyingWindows(bananas []int) map[Window]int {
	// sequences[window] -> banana count
	sequences := make(map[Window]int)
	diff := calculateDifferences(bananas)

	// sliding window
	for i := 0; i < len(diff)-3; i++ {
		window := Window{}
		for j := 0; j < windowSize; j++ {
			window.digits[j] = diff[i+j]
		}

		if _, exists := sequences[window]; !exists {
			sequences[window] = bananas[i+4]
		}
	}
	return sequences
}

func part2(data []int) int {
	memo := make(map[Window]int)

	for _, initial := range data {
		bananas := getSoldBananas(initial)
		sequences := getBuyingWindows(bananas)

		for window, count := range sequences {
			memo[window] += count
		}
	}

	bananas := 0
	for _, value := range memo {
		if value > bananas {
			bananas = value
		}
	}

	fmt.Println("Part 2:", bananas)
	return bananas
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
