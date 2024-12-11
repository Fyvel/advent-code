package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatData(rows []string) []int {
	stonesStr := strings.Split(string(rows[0]), " ")
	stones := make([]int, len(stonesStr))

	for i := range stonesStr {
		value, _ := strconv.Atoi(stonesStr[i])
		stones[i] = value
	}

	return stones
}

func processStone(stone int) []int {
	if stone == 0 {
		// replace value with 1
		return []int{1}
	} else {
		strValue := strconv.Itoa(stone)
		// split stone in 2
		if len(strValue)%2 == 0 {
			middle := len(strValue) / 2
			leftValue, _ := strconv.Atoi(strValue[:middle])
			rightValue, _ := strconv.Atoi(strValue[middle:])

			return []int{leftValue, rightValue}
		} else {
			// multiply by 2024
			return []int{stone * 2024}
		}
	}
}

func part1(stones []int) int {
	depth := 25

	for i := 0; i < depth; i++ {
		// fmt.Println(depth)

		nextStones := []int{}
		for _, stone := range stones {
			nextStones = append(nextStones, processStone(stone)...)
		}
		stones = nextStones
	}

	fmt.Println(len(stones))
	return len(stones)
}

func calculateStoneSize(stone int, depth int, memo map[string]int, syncLock *sync.Mutex) int {
	if depth == 0 {
		return 1
	}

	key := fmt.Sprintf("%d_%d", stone, depth)

	syncLock.Lock()
	cached, exists := memo[key]
	syncLock.Unlock()

	if exists {
		return cached
	}

	nextStones := processStone(stone)
	localSum := 0
	for _, nextStone := range nextStones {
		localSum += calculateStoneSize(nextStone, depth-1, memo, syncLock)
	}

	syncLock.Lock()
	memo[key] = localSum
	syncLock.Unlock()
	return localSum
}

func part2(stones []int) int {
	depth := 75
	results := make(chan int, len(stones))

	var syncLock sync.Mutex
	var waitGroup sync.WaitGroup
	memo := make(map[string]int)

	for _, stone := range stones {

		waitGroup.Add(1)

		go func(current int) {
			defer waitGroup.Done()

			result := calculateStoneSize(current, depth, memo, &syncLock)
			results <- result
			// fmt.Printf("stone %d - count %d \n", i, result)
		}(stone)
	}

	waitGroup.Wait()

	sum := 0
	for i := 0; i < len(stones); i++ {
		sum += <-results
	}

	fmt.Println(sum)
	return sum
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
