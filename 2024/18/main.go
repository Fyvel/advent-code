package main

import (
	"container/heap"
	"fmt"
	"os"
	"path/filepath"
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

type Vector struct {
	x, y int
}

func formatData(rows []string) []Vector {
	walls := make([]Vector, len(rows))
	for i, row := range rows {
		position := strings.Split(row, ",")
		x, _ := strconv.Atoi(position[0])
		y, _ := strconv.Atoi(position[1])
		walls[i] = Vector{x, y}
	}
	return walls
}

type Cell string

type Step struct {
	position  Vector
	direction int
}

type Simulator struct {
	size          int
	grid          map[Vector]Cell
	walls         []Vector
	path          []Step
	start         Vector
	end           Vector
	score         int
	exploredPaths [][]Step
}

func (s *Simulator) createGrid() *Simulator {
	s.grid = make(map[Vector]Cell)
	for r := 0; r < s.size; r++ {
		for c := 0; c < s.size; c++ {
			v := Vector{c, r}
			s.grid[v] = "."
		}
	}
	return s
}

func (s *Simulator) addWalls(newWalls []Vector) *Simulator {
	s.walls = append(s.walls, newWalls...)
	return s
}

func (s *Simulator) addWall(newWall Vector) *Simulator {
	s.walls = append(s.walls, newWall)
	return s
}

func (s *Simulator) isWall(v Vector) bool {
	for _, wall := range s.walls {
		if wall == v {
			return true
		}
	}
	return false
}

func (s *Simulator) isPath(v Vector) bool {
	for _, path := range s.path {
		if path.position == v {
			return true
		}
	}
	return false
}

func (s *Simulator) renderGrid() {
	// clear screen
	if len(s.exploredPaths) <= 1 {
		fmt.Print("\033[H\033[2J")
		for i := 0; i < s.size+2; i++ {
			fmt.Println()
		}
	}
	fmt.Print("\033[H")

	red := "\033[31m"
	green := "\033[32m"
	orange := "\033[33m"
	reset := "\033[0m"

	directionSymbols := []string{"↓", "→", "↑", "←"}
	render := make([][]string, s.size)
	for r := range render {
		render[r] = make([]string, s.size)
		for c := range render[r] {
			v := Vector{c, r}
			if s.isWall(v) {
				render[r][c] = orange + "#" + reset
			} else if s.score != -1 && s.score != int(^uint(0)>>1) && s.isPath(v) {
				render[r][c] = green + "O" + reset
			} else {
				foundDirection := false
				for _, path := range s.exploredPaths {
					for _, step := range path {
						if step.position == v {
							render[r][c] = red + directionSymbols[step.direction] + reset
							foundDirection = true
							break
						}
					}
					if foundDirection {
						break
					}
				}
				if !foundDirection {
					if v == s.start {
						render[r][c] = "S"
					} else if v == s.end {
						render[r][c] = "E"
					} else {
						render[r][c] = "."
					}
				}
			}
		}
		fmt.Print("\r" + strings.Join(render[r], " ") + "\n")
	}
	fmt.Println()
	time.Sleep(50 * time.Millisecond)
}

type QueueItem struct {
	position  Vector
	direction int
	score     int
	path      []Step
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
func (queue Queue) Less(a, b int) bool { return queue[a].score < queue[b].score }

func (s *Simulator) solve(render bool) *Simulator {
	// down right up left
	directions := []Vector{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	visited := make(map[string]bool)

	queue := &Queue{}
	heap.Init(queue)

	heap.Push(queue, &QueueItem{
		position: s.start,
		score:    0,
		path:     []Step{{s.start, 0}},
	})

	for queue.Len() > 0 {
		current := heap.Pop(queue).(*QueueItem)

		if current.score > s.score {
			continue
		}

		if current.position == s.end {
			s.score = current.score
			s.path = current.path
			return s
		}

		key := fmt.Sprintf("%d_%d_%d", current.position.x, current.position.y, current.direction)
		if visited[key] {
			continue
		}
		visited[key] = true

		for dirIndex, dir := range directions {
			next := Vector{
				x: current.position.x + dir.x,
				y: current.position.y + dir.y,
			}
			nextKey := fmt.Sprintf("%d_%d_%d", next.x, next.y, dirIndex)

			// out of bounds
			if next.x < 0 || next.x >= s.size ||
				next.y < 0 || next.y >= s.size ||
				s.isWall(next) || visited[nextKey] {
				continue
			}

			newPath := make([]Step, len(current.path))
			copy(newPath, current.path)
			newPath = append(newPath, Step{position: next, direction: dirIndex})

			// add to stack
			heap.Push(queue, &QueueItem{
				position: next,
				score:    current.score + 1,
				path:     newPath,
			})
		}

		s.path = current.path
		if render {
			s.exploredPaths = append(s.exploredPaths, current.path)
			s.renderGrid()
		}
	}

	s.renderGrid()
	s.score = -1
	return s
}

func part1(walls []Vector) int {
	mapSize := 71
	numberOfWalls := 1024
	renderSteps := false

	simulator := Simulator{
		size:  mapSize,
		walls: []Vector{},
		start: Vector{0, 0},
		end:   Vector{mapSize - 1, mapSize - 1},
		path:  []Step{},
		score: int(^uint(0) >> 1),
	}
	simulator.createGrid()
	simulator.addWalls(walls[:numberOfWalls])
	simulator.solve(renderSteps)
	simulator.renderGrid()

	fmt.Println("Shortest path:", simulator.score)
	return simulator.score
}

func part2(walls []Vector) string {
	mapSize := 71
	numberOfWalls := 1024
	renderSteps := false

	simulator := Simulator{
		size:  mapSize,
		walls: []Vector{},
		start: Vector{0, 0},
		end:   Vector{mapSize - 1, mapSize - 1},
		path:  []Step{},
		score: int(^uint(0) >> 1),
	}

	simulator.createGrid().addWalls(walls[:numberOfWalls]).solve(renderSteps)
	blockingWall := Vector{}
	for i := numberOfWalls; i < len(walls); i++ {
		wall := walls[i]

		simulation := Simulator{
			size:  mapSize,
			walls: make([]Vector, numberOfWalls),
			start: Vector{0, 0},
			end:   Vector{mapSize - 1, mapSize - 1},
			path:  []Step{},
			score: int(^uint(0) >> 1),
		}

		simulation.createGrid().addWalls(walls[:i]).addWall(wall).solve(renderSteps).renderGrid()

		if simulation.score == -1 {
			blockingWall = wall
			break
		}
	}

	key := fmt.Sprintf("%d,%d", blockingWall.x, blockingWall.y)
	fmt.Println("Blocking wall:", key)
	return key
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
