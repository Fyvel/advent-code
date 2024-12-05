package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FormattedData struct {
	rules   []string
	updates []string
}

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatData(rows []string) FormattedData {
	var rules, updates []string
	isEndOfRules := false

	for _, row := range rows {
		if row == "" {
			isEndOfRules = true
			continue
		}
		if !isEndOfRules {
			rules = append(rules, row)
		} else {
			updates = append(updates, row)
		}
	}
	return FormattedData{rules: rules, updates: updates}
}

func part1(data FormattedData) int {
	rulesSet := make(map[string]bool)
	for _, rule := range data.rules {
		rulesSet[rule] = true
	}

	var validUpdates [][]string

	for _, update := range data.updates {
		pages := strings.Split(update, ",")
		isValidUpdate := true

		for i := 0; i < len(pages)-1; i++ {
			failingRule := pages[i+1] + "|" + pages[i]
			if rulesSet[failingRule] {
				isValidUpdate = false
				break
			}
		}
		if isValidUpdate {
			validUpdates = append(validUpdates, pages)
		}
	}

	sum := 0
	for _, validUpdate := range validUpdates {
		mid, _ := strconv.Atoi(validUpdate[len(validUpdate)/2])
		sum += mid
	}

	fmt.Println(sum)
	return sum
}

func part2(data FormattedData) int {
	var correctedUpdates [][]int

	for _, update := range data.updates {
		pages := make([]int, 0)
		for _, p := range strings.Split(update, ",") {
			num, _ := strconv.Atoi(p)
			pages = append(pages, num)
		}

		rulesSet := make([][2]int, 0)
		for _, rule := range data.rules {
			parts := strings.Split(rule, "|")
			a, _ := strconv.Atoi(parts[0])
			b, _ := strconv.Atoi(parts[1])

			containsA := false
			containsB := false
			for _, p := range pages {
				containsA = containsA || p == a
				containsB = containsB || p == b
			}
			if containsA && containsB {
				rulesSet = append(rulesSet, [2]int{a, b})
			}
		}

		dependencies := make(map[int]int)
		for _, rule := range rulesSet {
			dependencies[rule[1]]++
		}

		corrected := make([]int, 0)
		for len(corrected) < len(pages) {
			for _, page := range pages {

				containsPage := false
				for _, c := range corrected {
					containsPage = containsPage || c == page
					if containsPage {
						break
					}
				}
				if containsPage {
					continue
				}

				if dependencies[page] < 1 {
					corrected = append(corrected, page)
					for _, rule := range rulesSet {
						if rule[0] == page {
							dependencies[rule[1]]--
						}
					}
				}
			}
		}

		original := strings.Trim(strings.Join(strings.Split(update, ","), ","), " ")
		correctedStr := ""
		for i, num := range corrected {
			if i > 0 {
				correctedStr += ","
			}
			correctedStr += strconv.Itoa(num)
		}

		if original != correctedStr {
			correctedUpdates = append(correctedUpdates, corrected)
		}
	}

	sum := 0
	for _, correctedUpdate := range correctedUpdates {
		sum += correctedUpdate[len(correctedUpdate)/2]
	}

	fmt.Println(sum)
	return sum
}

func main() {
	reports, err := readData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formattedData := formatData(reports)
	part1(formattedData)
	part2(formattedData)
}
