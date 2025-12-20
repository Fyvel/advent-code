package main

import (
	"aoc2025/day3/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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
	fmt.Println("Part 1:", sum)
}

func part2(banks []string) {
	sum := 0
	length := 12

	for _, bank := range banks {
		if len(bank) < length {
			continue
		}

		joltageNum, _ := processBank(bank, length, nil, 0)
		sum += joltageNum
	}
	fmt.Println("Part 2:", sum)
}

func processBank(bank string, length int, render *utils.Visualiser, bankIdx int) (int, map[int]bool) {
	joltage := ""
	startIdx := 0
	usedIndexes := make(map[int]bool)

	for idx := range length {
		remaining := length - idx - 1
		endIdx := len(bank) - remaining

		maxDigit := '0'
		maxDigitIdx := startIdx

		for i := startIdx; i < endIdx; i++ {
			if render != nil {
				render.UpdateSearching(bankIdx, bank, startIdx, i, maxDigitIdx, joltage)
			}
			if bank[i] > byte(maxDigit) {
				maxDigit = rune(bank[i])
				maxDigitIdx = i
			}
		}

		joltage += string(maxDigit)
		startIdx = maxDigitIdx + 1
		usedIndexes[maxDigitIdx] = true
	}

	joltageNum, _ := strconv.Atoi(joltage)
	if render != nil {
		render.Complete(bankIdx, bank, joltage, usedIndexes)
	}
	return joltageNum, usedIndexes
}

func part2visualAsync(banks []string) {
	length := 12
	render := utils.NewVisualiser(10*time.Millisecond, false)

	validBanks := []struct {
		idx  int
		bank string
	}{}
	for idx, bank := range banks {
		if len(bank) >= length {
			render.RegisterBank(idx)
			validBanks = append(validBanks, struct {
				idx  int
				bank string
			}{idx, bank})
		}
	}

	for range validBanks {
		fmt.Println()
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	sum := 0

	for _, vb := range validBanks {
		wg.Add(1)

		go func(bankIdx int, bank string) {
			defer wg.Done()

			joltageNum, _ := processBank(bank, length, render, bankIdx)

			mu.Lock()
			sum += joltageNum
			mu.Unlock()
		}(vb.idx, vb.bank)
	}

	wg.Wait()
	fmt.Println("\nSum of max values:", sum)
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

	withVisual := os.Getenv("AOC_VISUAL") == "1"
	if withVisual {
		part2visualAsync(formattedData)
	}
}
