package main

import (
	"container/heap"
	"fmt"
	"os"
	"path/filepath"
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
	grid      [][]string
	walls     map[Vector]bool
	start     Vector
	end       Vector
	score     int
	bestPaths map[Vector]bool
}

func formatData(rows []string) Maze {
	maze := Maze{
		grid:      make([][]string, len(rows)),
		walls:     make(map[Vector]bool),
		score:     int(^uint(0) >> 1),
		bestPaths: make(map[Vector]bool),
	}
	for r, row := range rows {
		maze.grid[r] = make([]string, len(row))
		for c, char := range row {
			v := Vector{c, r}
			switch char {
			case '#':
				maze.walls[v] = true
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
func (queue Queue) Less(a, b int) bool { return queue[a].score < queue[b].score }

func (maze *Maze) solve() int {
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

		if current.score > maze.score {
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
				maze.bestPaths = make(map[Vector]bool)
				maze.score = current.score
			}
			// mark best path
			for _, tile := range current.path {
				maze.bestPaths[tile] = true
			}
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
			score:     current.score + 1000,
			path:      newPathLeft,
		})

		rightDir := (current.direction + 1) % 4
		newPathRight := make([]Vector, len(current.path))
		copy(newPathRight, current.path)
		heap.Push(queue, &QueueItem{
			position:  current.position,
			direction: rightDir,
			score:     current.score + 1000,
			path:      newPathRight,
		})
	}

	if maze.score == int(^uint(0)>>1) {
		maze.score = -1
	}
	return maze.score
}

func part1(maze Maze) int {
	maze.solve()
	fmt.Println(maze.score)
	return maze.score
}

func part2(maze Maze) int {
	maze.solve()
	fmt.Println(len(maze.bestPaths))
	return len(maze.bestPaths)
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
