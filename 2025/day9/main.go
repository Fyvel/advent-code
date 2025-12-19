package main

import (
	"aoc2025/utils"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	viewportWidth  = 100
	viewportHeight = 60
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
	U, V  Coord
	Area  int
	Valid bool
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

func (p *Pair) IsWithinRectangle(coord Coord) bool {
	minR := int(math.Min(float64(p.U.R), float64(p.V.R)))
	maxR := int(math.Max(float64(p.U.R), float64(p.V.R)))
	minC := int(math.Min(float64(p.U.C), float64(p.V.C)))
	maxC := int(math.Max(float64(p.U.C), float64(p.V.C)))

	return coord.R >= minR && coord.R <= maxR && coord.C >= minC && coord.C <= maxC
}

func normaliseAndScale(coords Coords, width, height int) Coords {
	if len(coords) == 0 {
		return coords
	}

	minR, maxR := coords[0].R, coords[0].R
	minC, maxC := coords[0].C, coords[0].C
	for _, coord := range coords[1:] {
		minR = min(minR, coord.R)
		maxR = max(maxR, coord.R)
		minC = min(minC, coord.C)
		maxC = max(maxC, coord.C)
	}

	rangeR := maxR - minR
	rangeC := maxC - minC
	if rangeR == 0 {
		rangeR = 1
	}
	if rangeC == 0 {
		rangeC = 1
	}

	scaled := make(Coords, len(coords))
	scaleR := float64(width-1) / float64(rangeR*2)
	scaleC := float64(height-1) / float64(rangeC*2)

	for i, coord := range coords {
		normR := coord.R - minR
		normC := coord.C - minC

		scaledR := int(math.Round(float64(normR) * scaleR))
		scaledC := int(math.Round(float64(normC) * scaleC))

		if scaledR < 0 {
			scaledR = 0
		} else if scaledR >= width {
			scaledR = width - 1
		}
		if scaledC < 0 {
			scaledC = 0
		} else if scaledC >= height {
			scaledC = height - 1
		}

		scaled[i] = Coord{R: scaledR, C: scaledC}
	}

	return scaled
}

func (coords *Coords) BuildGrid() [][]string {
	var maxX, maxY int

	for _, coord := range *coords {
		if coord.R > maxX {
			maxX = coord.R
		}
		if coord.C > maxY {
			maxY = coord.C
		}
	}

	grid := make([][]string, maxY+1)
	for r := range grid {
		grid[r] = make([]string, maxX+1)
		for c := range grid[r] {
			grid[r][c] = "."
		}
	}

	for _, coord := range *coords {
		grid[coord.C][coord.R] = "#"
	}

	return grid
}

func part1(data Coords, withVisuals bool) {
	coords := make(Coords, len(data))
	copy(coords, data)

	displayCoords := normaliseAndScale(coords, viewportWidth, viewportHeight)

	var pairs []Pair
	for r := 0; r < len(coords)-1; r++ {
		for c := r + 1; c < len(coords); c++ {
			u := coords[r]
			v := coords[c]

			newPair := Pair{
				U:    u,
				V:    v,
				Area: calculateArea(u, v),
			}
			pairs = append(pairs, newPair)

			// visuals with scaled coordinates
			if !withVisuals || (r%100 != 0 && c%100 != 0) {
				continue
			}

			dU := displayCoords[r]
			dV := displayCoords[c]
			minR := int(math.Min(float64(dU.R), float64(dV.R)))
			maxR := int(math.Max(float64(dU.R), float64(dV.R)))
			minC := int(math.Min(float64(dU.C), float64(dV.C)))
			maxC := int(math.Max(float64(dU.C), float64(dV.C)))

			grid := displayCoords.BuildGrid()
			for gr := range grid {
				for gc := range grid[gr] {
					if gc >= minR && gc <= maxR && gr >= minC && gr <= maxC {
						if grid[gr][gc] == "." {
							grid[gr][gc] = "O"
						}
					}
				}
			}
			grid[dU.C][dU.R] = "U"
			grid[dV.C][dV.R] = "V"
			utils.RenderGrid(grid, r, c, nil, cellRenderer)
			time.Sleep(1 * time.Millisecond)
		}
	}

	sort.Slice(pairs, func(a, b int) bool { return pairs[a].Area > pairs[b].Area })

	largestRectangle := pairs[0]
	fmt.Println("Part 1:", largestRectangle.Area)
}

var cellRenderer = func(ctx utils.CellRenderContext) string {
	if ctx.Cell == "#" {
		return utils.BgRed + utils.White + ctx.Cell + utils.Reset
	}
	if ctx.Cell == "O" {
		return utils.BgGreen + utils.White + ctx.Cell + utils.Reset
	}
	if ctx.Cell == "U" || ctx.Cell == "V" {
		return utils.BgGreen + utils.Red + "#" + utils.Reset
	}
	if ctx.Cell == "0" {
		return utils.BgGreen + utils.White + "▧" + utils.Reset
	}
	return utils.BgBlack + utils.White + ctx.Cell + utils.Reset
}

func part2(data Coords, withVisuals bool) {
	redTiles := map[Coord]bool{}
	redTilesX := map[int][]Coord{}
	redTilesY := map[int][]Coord{}

	for _, coord := range data {
		redTiles[coord] = true
		redTilesX[coord.R] = append(redTilesX[coord.R], coord)
		redTilesY[coord.C] = append(redTilesY[coord.C], coord)
	}

	markedTiles := map[Coord]bool{}
	markedTilesY := map[int][]Coord{}

	for rt := range redTiles {
		for _, tile := range redTilesX[rt.R] {
			if rt == tile {
				continue
			}
			minY := min(rt.C, tile.C)
			maxY := max(rt.C, tile.C)

			for c := minY; c <= maxY; c++ {
				cell := Coord{R: rt.R, C: c}
				markedTiles[cell] = true
				markedTilesY[c] = append(markedTilesY[c], cell)
			}
		}

		for _, tile := range redTilesY[rt.C] {
			if rt == tile {
				continue
			}

			minX := min(rt.R, tile.R)
			maxX := max(rt.R, tile.R)

			for c := minX; c <= maxX; c++ {
				cell := Coord{R: c, C: rt.C}
				markedTiles[cell] = true
				markedTilesY[rt.C] = append(markedTilesY[rt.C], cell)
			}
		}
	}

	segments := map[int]Segment{}

	for c, coords := range markedTilesY {
		if len(coords) < 2 {
			continue
		}

		minX := math.MaxInt
		maxX := -1

		for _, coord := range coords {
			minX = min(minX, coord.R)
			maxX = max(maxX, coord.R)
		}

		segments[c] = Segment{Min: minX, Max: maxX}
	}

	largestRectangle := 0
	var largestU, largestV Coord

	displayCoords := normaliseAndScale(data, viewportWidth, viewportHeight)
	coordToDisplay := make(map[Coord]Coord)
	for i, coord := range data {
		coordToDisplay[coord] = displayCoords[i]
	}

	for c1 := range redTiles {
		for c2 := range redTiles {
			if c1 == c2 {
				continue
			}

			rectangle := calculateArea(c1, c2)

			if !checkTiles(c1, c2, segments) {
				continue
			}

			if withVisuals {
				renderRectangleVisuals(displayCoords, data, c1, c2, largestRectangle, coordToDisplay, largestU, largestV)
			}

			if rectangle <= largestRectangle {
				continue
			}

			largestRectangle = rectangle
			largestU = c1
			largestV = c2
		}
	}

	fmt.Println("Part 2:", largestRectangle)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData, true)
	part2(formattedData, true)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Segment struct {
	Min, Max int
}

func checkTiles(u, v Coord, segments map[int]Segment) bool {
	minX := min(u.R, v.R)
	maxX := max(u.R, v.R)
	minY := min(u.C, v.C)
	maxY := max(u.C, v.C)

	for y := minY; y <= maxY; y++ {
		area := segments[y]

		if minX < area.Min || maxX > area.Max {
			return false
		}
	}

	return true
}

func renderRectangleVisuals(displayCoords Coords, data Coords, c1 Coord, c2 Coord, largestRectangle int, coordToDisplay map[Coord]Coord, largestU Coord, largestV Coord) {
	dU := displayCoords[0]
	dV := displayCoords[0]
	for i, coord := range data {
		if coord == c1 {
			dU = displayCoords[i]
		}
		if coord == c2 {
			dV = displayCoords[i]
		}
	}

	minR := min(dU.R, dV.R)
	maxR := max(dU.R, dV.R)
	minC := min(dU.C, dV.C)
	maxC := max(dU.C, dV.C)
	grid := displayCoords.BuildGrid()

	var largestMinR, largestMaxR, largestMinC, largestMaxC int
	if largestRectangle > 0 {
		lU := coordToDisplay[largestU]
		lV := coordToDisplay[largestV]
		largestMinR = min(lU.R, lV.R)
		largestMaxR = max(lU.R, lV.R)
		largestMinC = min(lU.C, lV.C)
		largestMaxC = max(lU.C, lV.C)

		for gr := range grid {
			for gc := range grid[gr] {
				if gc >= largestMinR && gc <= largestMaxR && gr >= largestMinC && gr <= largestMaxC {
					if grid[gr][gc] == "." {
						grid[gr][gc] = "0"
					}
				}
			}
		}
	}

	for gr := range grid {
		for gc := range grid[gr] {
			if gc >= minR && gc <= maxR && gr >= minC && gr <= maxC {
				if grid[gr][gc] == "." {
					grid[gr][gc] = "O"
				}
			}
		}
	}
	grid[dU.C][dU.R] = "U"
	grid[dV.C][dV.R] = "V"

	utils.RenderGrid(grid, -1, -1, nil, cellRenderer)
	time.Sleep(1 * time.Millisecond)
}
