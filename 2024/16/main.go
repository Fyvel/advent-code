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
	grid  [][]string
	walls map[Vector]bool
	start Vector
	end   Vector
	score int
}

func formatData(rows []string) Maze {
	maze := Maze{
		grid:  make([][]string, len(rows)),
		walls: make(map[Vector]bool),
		score: int(^uint(0) >> 1),
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
	})

	for queue.Len() > 0 {
		current := queue.Pop().(*QueueItem)

		if current.score >= maze.score {
			continue
		}

		key := fmt.Sprintf("%d_%d_%d", current.position.x, current.position.y, current.direction)
		if _, exist := visited[key]; exist {
			if visited[key] <= current.score {
				continue
			}
		}
		visited[key] = current.score

		if current.position == maze.end {
			maze.score = current.score
			return current.score
		}

		// forward - left - right
		nextPosition := Vector{
			current.position.x + directions[current.direction].x,
			current.position.y + directions[current.direction].y,
		}
		if !maze.walls[nextPosition] && nextPosition.x >= 0 && nextPosition.y >= 0 &&
			nextPosition.y < len(maze.grid) && nextPosition.x < len(maze.grid[0]) {
			heap.Push(queue, &QueueItem{nextPosition, current.direction, current.score + 1})
		}

		leftDir := (current.direction + 3) % 4
		heap.Push(queue, &QueueItem{current.position, leftDir, current.score + 1000})

		rightDir := (current.direction + 1) % 4
		heap.Push(queue, &QueueItem{current.position, rightDir, current.score + 1000})
	}

	if maze.score == int(^uint(0)>>1) {
		maze.score = -1
		return -1
	}
	return maze.score
}

func part1(maze Maze) int {
	maze.solve()
	fmt.Println(maze.score)
	return maze.score

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
