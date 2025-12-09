package main

import (
	"aoc2025/utils"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatData(rows []string) Coords {
	var coords Coords
	for _, row := range rows {
		var x, y int
		fmt.Sscanf(row, "%d,%d", &x, &y)
		coords = append(coords, Coord{R: x, C: y})
	}
	return coords
}

type Coord struct {
	R, C int
}

func (a Coord) DistanceSq(b Coord) float64 {
	dx := float64(a.R - b.R)
	dy := float64(a.C - b.C)
	return dx*dx + dy*dy
}

func (a Coord) Distance(b Coord) float64 {
	return math.Sqrt(a.DistanceSq(b))
}

type Coords []Coord

type Pair struct {
	U, V Coord
	Area int
}

func calculateArea(u, v Coord) int {
	minR := int(math.Min(float64(u.R), float64(v.R)))
	maxR := int(math.Max(float64(u.R), float64(v.R)))
	minC := int(math.Min(float64(u.C), float64(v.C)))
	maxC := int(math.Max(float64(u.C), float64(v.C)))

	width := maxR - minR + 1
	height := maxC - minC + 1
	return height * width
}

// func (coords *Coords) BuildGrid() [][]string {
// 	var maxX, maxY int
// 	var padding = 1

// 	for _, coord := range *coords {
// 		if coord.R > maxX {
// 			maxX = coord.R
// 		}
// 		if coord.C > maxY {
// 			maxY = coord.C
// 		}
// 	}

// 	grid := make([][]string, maxY+1+padding)
// 	for r := range grid {
// 		grid[r] = make([]string, maxX+1+padding)
// 		for c := range grid[r] {
// 			grid[r][c] = "."
// 		}
// 	}

// 	for _, coord := range *coords {
// 		grid[coord.C][coord.R] = "#"
// 	}

// 	return grid
// }

// func fillArea(grid [][]string, u, v Coord, fill string) [][]string {
// 	minR := int(math.Min(float64(u.R), float64(v.R)))
// 	maxR := int(math.Max(float64(u.R), float64(v.R)))
// 	minC := int(math.Min(float64(u.C), float64(v.C)))
// 	maxC := int(math.Max(float64(u.C), float64(v.C)))

// 	for r := minC; r <= maxC; r++ {
// 		for c := minR; c <= maxR; c++ {
// 			if !(r == u.C && c == u.R) && !(r == v.C && c == v.R) {
// 				grid[r][c] = fill
// 			}
// 		}
// 	}

// 	return grid
// }

func part1(data Coords) {
	// grid := data.BuildGrid()
	// cellRenderer := func(ctx utils.CellRenderContext) string {
	// 	if ctx.Cell == "#" {
	// 		return "🔴" + utils.Reset
	// 	}
	// 	if ctx.Cell == "O" {
	// 		return "🟢" + utils.Reset
	// 	}
	// 	return "⬛" + utils.Reset
	// }

	var pairs []Pair
	for r := 0; r < len(data)-1; r++ {
		for c := r + 1; c < len(data); c++ {
			u := data[r]
			v := data[c]
			pairs = append(pairs, Pair{
				U:    u,
				V:    v,
				Area: calculateArea(u, v),
			})

			// tempGrid := make([][]string, len(grid))
			// for i := range grid {
			// 	tempGrid[i] = make([]string, len(grid[i]))
			// 	copy(tempGrid[i], grid[i])
			// }
			// tempGrid = fillArea(tempGrid, u, v, "O")

			// utils.EnterVisualMode()

			// utils.RenderGrid(tempGrid, r, c, nil, cellRenderer)
			// time.Sleep(55 * time.Millisecond)
		}
	}

	utils.ExitVisualMode()

	sort.Slice(pairs, func(a, b int) bool { return pairs[a].Area > pairs[b].Area })

	largestRectangle := pairs[0]
	// tempGrid := make([][]string, len(grid))
	// for i := range grid {
	// 	tempGrid[i] = make([]string, len(grid[i]))
	// 	copy(tempGrid[i], grid[i])
	// }
	// tempGrid = fillArea(tempGrid, largestRectangle.U, largestRectangle.V, "O")

	// utils.RenderGrid(tempGrid, -1, -1, nil, cellRenderer)

	fmt.Println("Part 1:", largestRectangle.Area)
}

func part2(data Coords) {
	fmt.Println("Part 2:", data)
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
