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

func formatData(rows []string) string {
	return rows[0]
}

func part1(diskMap string) int64 {
	var fileBlocks []string
	blockId := 0

	for i := 0; i < len(diskMap); i++ {
		char := diskMap[i]
		count := int(char - '0')
		for j := 0; j < count; j++ {
			if i%2 == 0 {
				fileBlocks = append(fileBlocks, fmt.Sprintf("%d", blockId))
			} else {
				fileBlocks = append(fileBlocks, ".")
			}
		}
		if i%2 == 0 {
			blockId++
		}
	}

	leftIndex := 0
	rightIndex := len(fileBlocks) - 1

	for leftIndex < rightIndex {
		left := fileBlocks[leftIndex]
		right := fileBlocks[rightIndex]

		if left == "." && right != "." {
			temp := fileBlocks[leftIndex]
			fileBlocks[leftIndex] = fileBlocks[rightIndex]
			fileBlocks[rightIndex] = temp

			leftIndex++
			rightIndex--
			continue
		}

		if right == "." {
			rightIndex--
		}
		if left != "." {
			leftIndex++
		}
	}

	var checksum int64
	for i := 0; i < len(fileBlocks); i++ {
		if fileBlocks[i] == "." {
			break
		}
		val, _ := strconv.Atoi(fileBlocks[i])
		checksum += int64(i) * int64(val)
	}

	// fmt.Println(strings.Join(fileBlocks, ""))
	fmt.Println(checksum)
	return checksum
}

type Block struct {
	value string
	start int
	size  int
}

func slidingWindow(fileBlocks []string) []string {
	result := make([]string, len(fileBlocks))
	copy(result, fileBlocks)

	var blocks []Block
	currentValue := result[0]
	currentStart := 0
	currentSize := 1

	// find blocks
	for i := 1; i < len(result); i++ {
		if result[i] == currentValue {
			currentSize++
		} else {
			blocks = append(blocks, Block{currentValue, currentStart, currentSize})
			currentValue = result[i]
			currentStart = i
			currentSize = 1
		}
	}
	blocks = append(blocks, Block{currentValue, currentStart, currentSize})

	// right to left
	for i := len(blocks) - 1; i >= 0; i-- {

		block := blocks[i]
		if block.value == "." {
			continue
		}

		var leftMostFreeSpace Block
		freeSpaceIndex := -1

		for j := 0; j < i; j++ {
			localBlock := blocks[j]
			if localBlock.value != "." {
				continue
			}
			if localBlock.size >= block.size {
				leftMostFreeSpace = localBlock
				freeSpaceIndex = j
				break
			}
		}

		if freeSpaceIndex != -1 {
			copy(result[leftMostFreeSpace.start:], result[block.start:block.start+block.size])

			for k := block.start; k < block.start+block.size; k++ {
				result[k] = "."
			}
			blocks[freeSpaceIndex].size -= block.size
			blocks[freeSpaceIndex].start += block.size
		}
	}

	return result
}

func part2(diskMap string) int64 {
	var fileBlocks []string
	blockId := 0

	for i := 0; i < len(diskMap); i++ {
		char := diskMap[i]
		count := int(char - '0')
		for j := 0; j < count; j++ {
			if i%2 == 0 {
				fileBlocks = append(fileBlocks, fmt.Sprintf("%d", blockId))
			} else {
				fileBlocks = append(fileBlocks, ".")
			}
		}
		if i%2 == 0 {
			blockId++
		}
	}

	// fmt.Println(strings.Join(fileBlocks, ""))

	fileBlocks = slidingWindow(fileBlocks)

	// fmt.Println(strings.Join(fileBlocks, ""))

	var checksum int64
	for i := 0; i < len(fileBlocks); i++ {
		if fileBlocks[i] == "." {
			continue
		}
		val, _ := strconv.Atoi(fileBlocks[i])
		checksum += int64(i) * int64(val)
	}

	fmt.Println(checksum)
	return checksum
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
