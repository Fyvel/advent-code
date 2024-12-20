package main

import (
	"container/heap"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

type Vector struct {
	x, y int
}
type Maze struct {
	grid        [][]string
	walls       map[Vector]bool
	start       Vector
	end         Vector
	score       int
	bestPaths   []Vector
	cheats      map[Vector]int // walls surrounded by at least 2 empty spaces
	cheatGroups map[int]int
}

func formatData(rows []string) Maze {
	maze := Maze{
		grid:        make([][]string, len(rows)),
		walls:       make(map[Vector]bool),
		score:       int(^uint(0) >> 1),
		bestPaths:   []Vector{},
		cheats:      make(map[Vector]int),
		cheatGroups: make(map[int]int),
	}
	// North, East, South, West
	directions := []Vector{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	for r, row := range rows {
		maze.grid[r] = make([]string, len(row))
		for c, char := range row {
			v := Vector{c, r}
			switch char {
			case '#':
				maze.walls[v] = true
				// check if wall is surrounded by at least 3 empty spaces in each direction
				emptySpaces := 0
				for _, direction := range directions {
					adjacent := Vector{v.x + direction.x, v.y + direction.y}
					if adjacent.x >= 0 && adjacent.x < len(row) &&
						adjacent.y >= 0 && adjacent.y < len(rows) {
						if _, exist := maze.walls[adjacent]; !exist {
							emptySpaces++
						}
					}
				}
				if emptySpaces >= 3 {
					maze.cheats[v] = emptySpaces
				}
			case 'S':
				maze.start = v
			case 'E':
				maze.end = v
			}
			maze.grid[r][c] = string(char)
		}
	}

	return maze
}

type QueueItem struct {
	position  Vector
	direction int
	score     int
	path      []Vector
}
type Queue []*QueueItem

func (queue *Queue) Pop() interface{} {
	item := (*queue)[0]
	*queue = (*queue)[1:]
	return item
}
func (queue *Queue) Push(item interface{}) { *queue = append(*queue, item.(*QueueItem)) }
func (queue Queue) Swap(a, b int)          { queue[a], queue[b] = queue[b], queue[a] }
func (queue Queue) Len() int               { return len(queue) }

// sort queue by lowest score
func (queue Queue) Less(a, b int) bool { return len(queue[a].path) < len(queue[b].path) }

func (maze *Maze) solve() *Maze {
	MINIMUM_IMPROVEMENT := 100
	// North, East, South, West
	directions := []Vector{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	visited := make(map[string]int)

	queue := &Queue{}
	heap.Init(queue)

	heap.Push(queue, &QueueItem{
		position:  maze.start,
		direction: 1,
		score:     0,
		path:      []Vector{maze.start},
	})

	for queue.Len() > 0 {
		current := queue.Pop().(*QueueItem)

		if maze.score != int(^uint(0)>>1) &&
			len(current.path) > len(maze.bestPaths)-MINIMUM_IMPROVEMENT {
			continue
		}

		if current.score >= maze.score {
			continue
		}

		key := fmt.Sprintf("%d_%d_%d", current.position.x, current.position.y, current.direction)
		if _, exist := visited[key]; exist {
			if visited[key] < current.score {
				continue
			}
		}
		visited[key] = current.score

		if current.position == maze.end {
			// reset for new best score
			if current.score < maze.score {
				maze.bestPaths = []Vector{}
				maze.score = current.score
			}
			// mark best path
			for _, tile := range current.path {
				if tile != maze.start {
					maze.bestPaths = append(maze.bestPaths, tile)
				}
			}
			fmt.Println("New best path:", len(maze.bestPaths))
			continue
		}

		// forward - left - right
		nextPosition := Vector{
			current.position.x + directions[current.direction].x,
			current.position.y + directions[current.direction].y,
		}
		if !maze.walls[nextPosition] && nextPosition.x >= 0 && nextPosition.y >= 0 &&
			nextPosition.y < len(maze.grid) && nextPosition.x < len(maze.grid[0]) {
			newPath := make([]Vector, len(current.path))
			copy(newPath, current.path)
			newPath = append(newPath, nextPosition)
			heap.Push(queue, &QueueItem{
				position:  nextPosition,
				direction: current.direction,
				score:     current.score + 1,
				path:      newPath,
			})
		}

		leftDir := (current.direction + 3) % 4
		newPathLeft := make([]Vector, len(current.path))
		copy(newPathLeft, current.path)
		heap.Push(queue, &QueueItem{
			position:  current.position,
			direction: leftDir,
			score:     current.score + 1,
			path:      newPathLeft,
		})

		rightDir := (current.direction + 1) % 4
		newPathRight := make([]Vector, len(current.path))
		copy(newPathRight, current.path)
		heap.Push(queue, &QueueItem{
			position:  current.position,
			direction: rightDir,
			score:     current.score + 1,
			path:      newPathRight,
		})
	}

	if maze.score == int(^uint(0)>>1) {
		maze.score = -1
	}

	return maze
}

func (m *Maze) copy() *Maze {
	newMaze := Maze{
		grid:      make([][]string, len(m.grid)),
		walls:     make(map[Vector]bool),
		start:     m.start,
		end:       m.end,
		score:     m.score,
		bestPaths: make([]Vector, len(m.bestPaths)),
		cheats:    make(map[Vector]int),
	}

	for i := range m.grid {
		newMaze.grid[i] = make([]string, len(m.grid[i]))
		copy(newMaze.grid[i], m.grid[i])
	}
	for k, v := range m.walls {
		newMaze.walls[k] = v
	}
	copy(newMaze.bestPaths, m.bestPaths)
	for k, v := range m.cheats {
		newMaze.cheats[k] = v
	}
	return &newMaze
}

func (m *Maze) isPath(v Vector) bool {
	for _, path := range m.bestPaths {
		if path == v {
			return true
		}
	}
	return false
}

func (m *Maze) renderGrid() {
	// clear screen
	// fmt.Print("\033[H\033[2J")
	fmt.Println()

	red := "\033[31m"
	green := "\033[32m"
	orange := "\033[33m"
	reset := "\033[0m"
	fmt.Println()

	// directionSymbols := []string{"↓", "→", "↑", "←"}
	render := make([][]string, len(m.grid))
	for r := range render {
		render[r] = make([]string, len(m.grid[r]))
		for c := range render[r] {
			v := Vector{c, r}
			if m.walls[v] {
				render[r][c] = orange + "#" + reset
			} else if m.score != -1 && m.score != int(^uint(0)>>1) && m.isPath(v) {
				render[r][c] = green + "O" + reset
			} else {
				if v == m.start {
					render[r][c] = "S"
				} else if v == m.end {
					render[r][c] = "E"
				} else {
					render[r][c] = "."
				}
			}
			if m.cheats[v] > 0 {
				render[r][c] = red + "#" + reset
			}
		}
		fmt.Print("\r" + strings.Join(render[r], " ") + "\n")
	}
	fmt.Println()
}

func (maze *Maze) activateCheat(v Vector) *Maze {
	maze.walls[v] = false
	return maze
}

func (m *Maze) groupCheatsByLength(cheatsShortestPaths map[Vector]int) *Maze {
	// group cheats by path length (keep only gaining at least 100 picoseconds)
	cheatsByPathLength := make(map[int]int)
	for _, pathLength := range cheatsShortestPaths {
		gain := len(m.bestPaths) - pathLength
		cheatsByPathLength[gain]++
	}
	m.cheatGroups = cheatsByPathLength
	return m
}

func part1(input Maze) {
	maze := input.copy().solve()
	maze.renderGrid()
	fmt.Println("Shortest path:", len(maze.bestPaths))
	fmt.Println("Potential cheats:", len(input.cheats))

	results := make(chan struct {
		cheat  Vector
		length int
	})
	workerPool := make(chan struct{}, runtime.NumCPU())

	for cheat := range maze.cheats {
		go func(c Vector) {
			workerPool <- struct{}{}
			cheatMaze := maze.copy()
			cheatMaze.activateCheat(c)
			cheatMaze.solve()
			results <- struct {
				cheat  Vector
				length int
			}{c, len(cheatMaze.bestPaths)}
			<-workerPool
		}(cheat)
	}

	cheatsShortestPaths := make(map[Vector]int)
	totalCheats := len(maze.cheats)
	processed := 0

	for result := range results {
		if result.length <= len(maze.bestPaths) {
			cheatsShortestPaths[result.cheat] = result.length
		}
		processed++
		if processed == totalCheats {
			close(results)
			break
		}
	}

	maze.groupCheatsByLength(cheatsShortestPaths)

	fmt.Println("Cheats by path length:")
	gains := make([]int, 0, len(maze.cheatGroups))
	for gain := range maze.cheatGroups {
		gains = append(gains, gain)
	}
	sort.Ints(gains)
	for _, gain := range gains {
		fmt.Printf("There are %d cheats that save %d picoseconds\n", maze.cheatGroups[gain], gain)
	}

	fmt.Println()
}

// func part2() {}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData)
	// part2(formattedData)
}
