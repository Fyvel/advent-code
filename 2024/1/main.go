package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func readData(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func formatData(rows []string) ([]int, []int, error) {
	var left []int
	var right []int
	for _, row := range rows {
		parts := strings.Split(row, "   ")
		if len(parts) == 2 {
			a, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, nil, err
			}
			b, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, nil, err
			}
			left = append(left, a)
			right = append(right, b)
		}
	}
	return left, right, nil
}

func part1(left, right []int) int {
	sort.Ints(left)
	sort.Ints(right)
	distance := 0
	for i := 0; i < len(left); i++ {
		distance += abs(left[i] - right[i])
	}
	fmt.Printf("{ distance: %d }\n", distance)
	return distance
}

func part2(left, right []int) int {
	similarity := 0
	mapRight := make(map[int]int)
	for _, value := range right {
		mapRight[value]++
	}
	for _, value := range left {
		similarity += value * mapRight[value]
	}
	fmt.Printf("{ similarity: %d }\n", similarity)
	return similarity
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	filePath := filepath.Join(".", "data.txt")
	rows, err := readData(filePath)
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	left, right, err := formatData(rows)
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	part1(left, right)
	part2(left, right)
}
