package main

import (
	"aoc2025/utils"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatData(rows []string) []Vector {
	var vectors []Vector
	for _, row := range rows {
		parts := strings.Split(row, ",")
		x, _ := strconv.ParseFloat(parts[0], 64)
		y, _ := strconv.ParseFloat(parts[1], 64)
		z, _ := strconv.ParseFloat(parts[2], 64)
		vectors = append(vectors, Vector{X: x, Y: y, Z: z})
	}
	return vectors
}

type Vector struct {
	X, Y, Z float64
}

func (a Vector) DistanceSq(b Vector) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return dx*dx + dy*dy + dz*dz
}

func (a Vector) Distance(b Vector) float64 {
	return math.Sqrt(a.DistanceSq(b))
}

type Pair struct {
	U, V     Vector
	Distance float64
}

type Circuits struct {
	counts  []int
	circuit []int
}

func NewCircuitsManager(size int) *Circuits {
	circuit := make([]int, size)
	group := make([]int, size)
	for i := range size {
		circuit[i] = i
		group[i] = 1
	}
	return &Circuits{
		counts:  group,
		circuit: circuit,
	}
}

func (c *Circuits) Search(junctionBox int) int {
	// base: junction box is its own group
	if c.circuit[junctionBox] == junctionBox {
		return junctionBox
	}
	// recursive: keep searching
	root := c.Search(c.circuit[junctionBox])
	c.circuit[junctionBox] = root
	return root
}

func (c *Circuits) Connect(u, v int) {
	circuitU := c.Search(u)
	circuitV := c.Search(v)
	if circuitU == circuitV {
		return // "already in the same circuit, nothing happens!"
	}
	// order by size
	if c.counts[circuitU] < c.counts[circuitV] {
		circuitU, circuitV = circuitV, circuitU // swap
	}

	c.circuit[circuitV] = circuitU
	c.counts[circuitU] += c.counts[circuitV]
}

func part1(junctionBoxes JunctionBoxes, withVisual bool) {

	pairs := junctionBoxes.buildPairs()
	jbIdxMap := junctionBoxes.buildIndexMap()

	connectionCount := 1000
	manager := NewCircuitsManager(len(junctionBoxes))

	for _, pair := range pairs {
		if connectionCount <= 0 {
			break
		}
		u := jbIdxMap[pair.U]
		v := jbIdxMap[pair.V]

		manager.Connect(u, v)
		connectionCount--
	}

	circuitsMap := make(map[int][]Vector)
	multiplyTop3 := 1

	for i := range junctionBoxes {
		cId := manager.Search(i)
		circuitsMap[cId] = append(circuitsMap[cId], junctionBoxes[i])
	}

	var circuits [][]Vector
	for _, circuit := range circuitsMap {
		circuits = append(circuits, circuit)
	}
	sort.Slice(circuits, func(a, b int) bool {
		return len(circuits[a]) > len(circuits[b])
	})

	if withVisual {
		top3 := 3
		renderCircuits(circuits, top3)

		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("%d circuits found\n", len(circuits)))

		for i := 0; i < top3 && i < len(circuits); i++ {
			multiplyTop3 *= len(circuits[i])
			buf.WriteString(fmt.Sprintf("Circuit #%d -> %d junction boxes\n", i+1, len(circuits[i])))
		}

		buf.WriteString(fmt.Sprintf("Part 1: %d\n", multiplyTop3))
		fmt.Print(buf.String())
		return
	}

	fmt.Println("Part 1:", multiplyTop3)
}

type JunctionBoxes []Vector

func (junctionBoxes JunctionBoxes) buildPairs() []Pair {
	var pairs []Pair
	for r := 0; r < len(junctionBoxes)-1; r++ {
		for c := r + 1; c < len(junctionBoxes); c++ {
			distSq := junctionBoxes[r].DistanceSq(junctionBoxes[c])
			pairs = append(pairs, Pair{
				U:        junctionBoxes[r],
				V:        junctionBoxes[c],
				Distance: math.Sqrt(distSq),
			})
		}
	}

	sort.Slice(pairs, func(a, b int) bool { return pairs[a].Distance < pairs[b].Distance })
	return pairs
}

func (junctionBoxes JunctionBoxes) buildIndexMap() map[Vector]int {
	jbIdxMap := make(map[Vector]int)
	for i, jb := range junctionBoxes {
		jbIdxMap[jb] = i
	}
	return jbIdxMap
}

func renderCircuits(circuits [][]Vector, depth int) {
	cellRenderer := func(ctx utils.CellRenderContext) string {
		if ctx.Cell == "#" {
			return utils.BgOrange + utils.White + ctx.Cell + utils.Reset
		}
		if ctx.Cell == "." {
			return utils.BgBlack + utils.White + ctx.Cell + utils.Reset
		}
		return ctx.Cell
	}

	for i := range depth {
		circuit := circuits[i]
		gridSize := 42
		grid := make([][]string, gridSize)
		for r := range grid {
			grid[r] = make([]string, gridSize)
			for c := range grid[r] {
				grid[r][c] = "."
			}
		}

		for _, jb := range circuit {
			x := int(jb.X) % gridSize
			y := int(jb.Y) % gridSize
			grid[y][x] = "#"
			utils.RenderGrid(grid, y, x, nil, cellRenderer)
		}
		time.Sleep(55 * time.Millisecond)
	}
}

func part2(junctionBoxes JunctionBoxes, withVisual bool) {

	pairs := junctionBoxes.buildPairs()
	jbIdxMap := junctionBoxes.buildIndexMap()

	manager := NewCircuitsManager(len(junctionBoxes))

	multiplyLastConnectionX := 1
	for _, pair := range pairs {
		u := jbIdxMap[pair.U]
		v := jbIdxMap[pair.V]

		manager.Connect(u, v)

		if manager.counts[manager.Search(u)] == len(junctionBoxes) {
			multiplyLastConnectionX = int(pair.U.X) * int(pair.V.X)
			break
		}
	}

	circuitsMap := make(map[int][]Vector)

	for i := range junctionBoxes {
		cId := manager.Search(i)
		circuitsMap[cId] = append(circuitsMap[cId], junctionBoxes[i])
	}

	var maxCircuitPath []Vector
	for _, circuit := range circuitsMap {
		if len(circuit) > len(maxCircuitPath) {
			maxCircuitPath = circuit
		}
	}

	if withVisual {
		renderCircuits([][]Vector{maxCircuitPath}, 1)
	}

	fmt.Printf("Part 2: %d\n", multiplyLastConnectionX)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	withVisual := os.Getenv("AOC_VISUAL") == "1"

	defer fmt.Print(utils.ShowCursor)

	formattedData := formatData(data)
	part1(formattedData, withVisual)

	if withVisual {
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Print(utils.ShowCursor)
	part2(formattedData, withVisual)
}
