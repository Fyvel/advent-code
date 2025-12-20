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
	return strings.Split(string(data), ","), nil
}

func formatData(rows []string) [][]int {
	intervals := make([][]int, len(rows))

	for i, row := range rows {
		intervalStr := strings.Split(row, "-")
		firstId, _ := strconv.Atoi(intervalStr[0])
		lastId, _ := strconv.Atoi(intervalStr[1])
		intervals[i] = getIdsInRange(firstId, lastId)
	}

	return intervals
}

func getIdsInRange(firstId int, lastId int) []int {
	allIds := []int{}

	for i := firstId; i <= lastId; i++ {
		allIds = append(allIds, i)
	}

	return allIds
}

func part1(data [][]int) {
	invalidIds := [][]int{}
	sum := 0

	for i, interval := range data {
		invalidIds = append(invalidIds, []int{})
		for _, id := range interval {
			idStr := strconv.Itoa(id)
			if len(idStr)%2 != 0 {
				continue
			}

			middle := len(idStr) / 2
			left := idStr[:middle]
			right := idStr[middle:]

			if left == right {
				invalidIds[i] = append(invalidIds[i], id)
				sum += id
			}
		}
	}

	fmt.Println("Part 1:", sum)
}

func part2(data [][]int) {
	invalidIds := [][]int{}
	sum := 0

	for i, interval := range data {
		invalidIds = append(invalidIds, []int{})
		for _, id := range interval {
			if hasRepeatedSequence(strconv.Itoa(id)) {
				invalidIds[i] = append(invalidIds[i], id)
				sum += id
			}
		}
	}

	fmt.Println("Part 2:", sum)
}

func hasRepeatedSequence(id string) bool {
	idTwice := id + id
	trimmed := idTwice[1 : len(idTwice)-1]
	hasId := strings.Contains(trimmed, id)
	return hasId
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
