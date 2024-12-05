package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func readData() ([]string, error) {
	file, err := os.Open(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func formatData(rows []string) [][]rune {
	var grid [][]rune
	for _, row := range rows {
		grid = append(grid, []rune(row))
	}
	return grid
}

func part1(grid [][]rune) int {
	count := 0
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, 1}, {0, -1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[0]); col++ {
			for _, direction := range directions {
				if dfs(grid, row, col, direction, "XMAS") {
					count++
				}
			}
		}
	}

	fmt.Println(count)
	return count
}

func dfs(grid [][]rune, row, col int, direction []int, word string) bool {
	r, c := row, col
	for _, letter := range word {
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) {
			return false
		}
		if grid[r][c] != letter {
			return false
		}
		r += direction[0]
		c += direction[1]
	}
	return true
}

func part2(grid [][]rune) int {
	count := 0
	directions := [][]int{
		{1, 1},  // right down
		{1, -1}, // left down
	}

	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[0]); col++ {
			if grid[row][col] != 'M' && grid[row][col] != 'S' {
				continue
			}

			word := "MAS"
			if grid[row][col] == 'S' {
				word = "SAM"
			}
			middle := len(word) / 2

			if col+len(word)-1 >= len(grid[0]) {
				continue
			}

			crossingWord := "MAS"
			if grid[row][col+len(word)-1] == 'S' {
				crossingWord = "SAM"
			}

			if dfs(grid, row, col, directions[0], word) &&
				dfs(grid, row, col+len(word)-1, directions[1], crossingWord) &&
				grid[row+middle][col+middle] == rune(word[middle]) {
				count++
			}
		}
	}

	fmt.Println(count)
	return count
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
