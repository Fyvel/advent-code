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
	return strings.Split(string(data), "\n\n"), nil
}

type State struct {
	Areas    []Area
	Presents []Present
}

type Present struct {
	shape [][]rune
	size  int
}

type Area struct {
	size          int
	grid          [][]rune
	presentCounts []int
}

func formatData(rows []string) State {
	var formatted State

	for _, row := range rows {

		// presents
		if !strings.Contains(row, "x") {
			presentLines := strings.Split(strings.TrimSpace(row), "\n")[1:]
			var present Present
			for _, line := range presentLines {
				presentRow := []rune(line)
				present.shape = append(present.shape, presentRow)
				present.size += strings.Count(line, "#")
			}
			formatted.Presents = append(formatted.Presents, present)
			continue
		}

		// areas
		if strings.Contains(row, "x") {
			parts := strings.Split(row, "\n")
			for _, part := range parts {
				subParts := strings.SplitN(part, ": ", 2)
				if len(subParts) != 2 {
					continue
				}

				var area Area
				areaParts := strings.Split(strings.TrimSpace(subParts[0]), "x")
				w, _ := strconv.Atoi(areaParts[0])
				l, _ := strconv.Atoi(areaParts[1])
				area.size = w * l
				area.grid = make([][]rune, l)
				for r := range l {
					area.grid[r] = make([]rune, w)
					for c := range w {
						area.grid[r][c] = '.'
					}
				}

				presentsStr := strings.Split(strings.TrimSpace(subParts[1]), " ")
				for _, ps := range presentsStr {
					p, _ := strconv.Atoi(ps)
					area.presentCounts = append(area.presentCounts, p)
				}
				formatted.Areas = append(formatted.Areas, area)
			}
		}
	}
	return formatted
}

func part1(data State) {
	sum := 0

	for _, area := range data.Areas {
		minPresentSize := 0
		for idx, pCount := range area.presentCounts {
			minPresentSize += data.Presents[idx].size * pCount
		}

		// possible fit
		if minPresentSize < area.size {
			sum++
		}
	}
	fmt.Println("Part 1:", sum)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData)
}
