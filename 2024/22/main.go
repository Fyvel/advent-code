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
		// multiply by 64 and mix & prune
		secret ^= (secret * 64)
		secret %= 16777216

		// divide by 32 and mix & prune
		secret ^= (secret / 32)
		secret %= 16777216

		// multiply by 2048 and mix & prune
		secret ^= (secret * 2048)
		secret %= 16777216
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

// func part2(data []int) { }

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
