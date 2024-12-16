package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

type Inputs struct {
	grid  [][]string
	moves []Vector
}

func formatData(rows []string) Inputs {
	grid := [][]string{}
	moves := []Vector{}

	parseGrid := true
	for _, row := range rows {
		if row == "" {
			parseGrid = false
			continue
		}
		if parseGrid {
			grid = append(grid, strings.Split(row, ""))
		}
		if !parseGrid {
			remainingRowsAsString := strings.Join(rows, "")
			for _, char := range remainingRowsAsString {
				switch char {
				case '^':
					moves = append(moves, Vector{0, -1})
				case 'v':
					moves = append(moves, Vector{0, 1})
				case '<':
					moves = append(moves, Vector{-1, 0})
				case '>':
					moves = append(moves, Vector{1, 0})
				}
			}
			break
		}
	}
	return Inputs{grid, moves}
}

type Vector struct {
	x, y int
}

type Entity struct {
	start Vector
	end   Vector
}

type Game struct {
	robot    Entity
	walls    map[Entity]bool
	boxes    map[Entity]bool
	grid     [][]string
	isScaled bool
}

func (g *Game) isWall(v Vector) bool {
	for wall := range g.walls {
		if v.x >= wall.start.x && v.x <= wall.end.x &&
			v.y >= wall.start.y && v.y <= wall.end.y {
			return true
		}
	}
	return false
}

func (g *Game) isBox(v Vector) bool {
	for box := range g.boxes {
		if v.y == box.start.y &&
			v.x >= box.start.x && v.x <= box.end.x {
			return true
		}
	}
	return false
}

func (g *Game) getBox(v Vector) Entity {
	for box := range g.boxes {
		if v.y == box.start.y && // Same row
			v.x >= box.start.x && v.x <= box.end.x { // Within box width
			return box
		}
	}
	return Entity{}
}

func (g *Game) render() [][]string {
	if g.isScaled {
		return g.renderScale()
	}

	fmt.Println()
	// create board
	board := make([][]string, len(g.grid))
	for r := range g.grid {
		board[r] = make([]string, len(g.grid[r]))
		for c := range board[r] {
			board[r][c] = "."
		}
	}

	// place walls
	for wall := range g.walls {
		board[wall.start.y][wall.start.x] = "#"
	}

	// place boxes
	for box := range g.boxes {
		board[box.start.y][box.start.x] = "O"
	}

	// place robot
	board[g.robot.start.y][g.robot.start.x] = "@"

	// print board
	for _, row := range board {
		fmt.Println(strings.Join(row, ""))
	}

	return board
}

func isForwardBox(box Entity) bool {
	return box.start.x < box.end.x
}

func (g *Game) renderScale() [][]string {
	fmt.Println()
	// create board
	board := make([][]string, len(g.grid))
	for r := range g.grid {
		board[r] = make([]string, len(g.grid[r])*2)
		for c := range board[r] {
			board[r][c] = "."
		}
	}

	// place walls
	for wall := range g.walls {
		r := wall.start.y
		board[r][wall.start.x] = "#"
		board[r][wall.end.x] = "#"
	}

	// place boxes
	for box := range g.boxes {
		if isForwardBox(box) {
			r := box.start.y
			board[r][box.start.x] = "["
			board[r][box.end.x] = "]"
		}
	}

	// place robot
	r := g.robot.start.y
	board[r][g.robot.start.x] = "@"

	for _, row := range board {
		fmt.Println(strings.Join(row, ""))
	}
	return board
}

func (g *Game) setup(grid [][]string) *Game {
	g.grid = grid
	g.walls = make(map[Entity]bool)
	g.boxes = make(map[Entity]bool)

	for r, row := range grid {
		for c, col := range row {
			baseX := c
			if g.isScaled {
				baseX = c * 2
			}

			entity := Entity{
				Vector{baseX, r},
				Vector{baseX, r},
			}

			switch col {
			case "#":
				if g.isScaled {
					entity.end = Vector{baseX + 1, r}
				}
				g.walls[entity] = true
			case "@":
				g.robot = entity
			case "O", "[":
				if g.isScaled {
					entity.end = Vector{baseX + 1, r}
					forward := Entity{Vector{baseX, r}, Vector{baseX + 1, r}}
					backward := Entity{Vector{baseX + 1, r}, Vector{baseX, r}}
					g.boxes[forward] = true
					g.boxes[backward] = true
				} else {
					g.boxes[entity] = true
				}
			}
		}
	}
	return g
}

func (g *Game) moveRobot(move Vector) bool {
	nextPosition := Entity{
		Vector{g.robot.start.x + move.x, g.robot.start.y + move.y},
		Vector{g.robot.end.x + move.x, g.robot.end.y + move.y},
	}

	if g.isWall(nextPosition.start) || g.isWall(nextPosition.end) {
		return false
	}

	if g.isBox(nextPosition.start) || g.isBox(nextPosition.end) {
		leftBox, rightBox := g.getBox(nextPosition.start), g.getBox(nextPosition.end)
		if !g.moveBox(leftBox, move, 0) || !g.moveBox(rightBox, move, 0) {
			return false
		}

		// check for backwards or forwards cell
		if move.x != 0 {
			checkPosition := nextPosition.start
			if move.x > 0 {
				checkPosition = nextPosition.end
			}
			if g.isBox(checkPosition) {
				adjacentBox := g.getBox(checkPosition)
				if !g.moveBox(adjacentBox, move, 0) {
					return false
				}
			}
		}

		if move.y != 0 {
			checkPosition := nextPosition.start
			if move.y > 0 {
				checkPosition = nextPosition.end
			}
			if g.isBox(checkPosition) {
				adjacentBox := g.getBox(checkPosition)
				if !g.moveBox(adjacentBox, move, 0) {
					return false
				}
			}
		}
	}

	g.robot = nextPosition
	return true
}

func (g *Game) moveBox(box Entity, direction Vector, depth int) bool {
	if depth > max(len(g.grid), len(g.grid[0])) {
		return false
	}

	nextPosition := Entity{
		Vector{box.start.x + direction.x, box.start.y + direction.y},
		Vector{box.end.x + direction.x, box.end.y + direction.y},
	}

	if g.isWall(nextPosition.start) || g.isWall(nextPosition.end) {
		return false
	}

	if direction.x != 0 {
		// check for backwards or forwards box
		checkPosition := nextPosition.start
		if direction.x > 0 {
			checkPosition = nextPosition.end
		}
		if g.isBox(checkPosition) {
			adjacentBox := g.getBox(checkPosition)
			if !g.moveBox(adjacentBox, direction, depth+1) {
				return false
			}
		}
	}

	if direction.y != 0 {
		// check up and down of the box
		if g.isBox(nextPosition.start) {
			adjacentBox := g.getBox(nextPosition.start)
			if !g.moveBox(adjacentBox, direction, depth+1) {
				return false
			}
		}
		if g.isBox(nextPosition.end) {
			adjacentBox := g.getBox(nextPosition.end)
			if !g.moveBox(adjacentBox, direction, depth+1) {
				return false
			}
		}
	}

	delete(g.boxes, box)
	delete(g.boxes, Entity{box.end, box.start})

	g.boxes[nextPosition] = true
	return true
}

func part1(inputs Inputs) {
	game := Game{}
	game.setup(inputs.grid)
	game.render()

	for _, move := range inputs.moves {
		game.moveRobot(move)
	}

	// for m, move := range inputs.moves {
	// 	game.render()
	// 	if game.moveRobot(move) {
	// 		fmt.Printf("Move: (%d,%d)\n", move.x, move.y)
	// 	} else {
	// 		fmt.Printf("Can't move: (%d,%d)\n", move.x, move.y)
	// 	}
	// 	fmt.Printf("Move: %d out of %d\n", m+1, len(inputs.moves))
	// 	time.Sleep(50 * time.Millisecond)
	// }

	game.render()

	sum := 0
	for box := range game.boxes {
		sum += (box.start.y*100 + box.start.x)
	}

	fmt.Println("Part 1:", sum)
}

func part2(inputs Inputs) {
	game := Game{isScaled: true}
	game.setup(inputs.grid)
	game.render()

	// for _, move := range inputs.moves {
	// 	game.moveRobot(move)
	// }

	for m, move := range inputs.moves {
		var moveStr string
		switch move {
		case Vector{0, -1}:
			moveStr = "^"
		case Vector{0, 1}:
			moveStr = "v"
		case Vector{-1, 0}:
			moveStr = "<"
		case Vector{1, 0}:
			moveStr = ">"
		}
		game.render()
		if game.moveRobot(move) {
			fmt.Printf("Move: %s\n", moveStr)
		} else {
			fmt.Printf("Move: %s - Nope\n", moveStr)
		}
		fmt.Printf("Move: %d out of %d\n", m+1, len(inputs.moves))
		time.Sleep(50 * time.Millisecond)
	}

	game.render()

	sum := 0
	for box := range game.boxes {
		// count only forwards boxes
		if isForwardBox(box) {
			sum += (box.start.y*100 + box.start.x)
		}
	}

	fmt.Println("Part 2:", sum)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	// part1(formattedData)
	part2(formattedData)
}

func main2() {
	data, err := readData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	formattedData := formatData(data)

	fmt.Println("Choose mode:")
	fmt.Println("1. Auto run (using data.txt)")
	fmt.Println("2. Creator mode (WASD to move, Q to quit)")

	var choice string
	fmt.Scanln(&choice)

	if choice == "1" {
		part1(formattedData)
		part2(formattedData)
	} else if choice == "2" {
		game := Game{isScaled: true}
		game.setup(formattedData.grid)

		fmt.Print("\033[H\033[2J")
		fmt.Println("Use WASD to move, Q to quit")

		for {
			game.render()
			move := getPlayerMove()
			game.moveRobot(move)
			fmt.Print("\033[H\033[2J")
		}
	}
}

func getPlayerMove() Vector {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	char, _, err := keyboard.GetKey()
	if err != nil {
		return Vector{0, 0}
	}

	switch char {
	case 'w', 'W':
		return Vector{0, -1}
	case 's', 'S':
		return Vector{0, 1}
	case 'a', 'A':
		return Vector{-1, 0}
	case 'd', 'D':
		return Vector{1, 0}
	case 'q', 'Q':
		os.Exit(1)
		return Vector{0, 0}
	default:
		return Vector{0, 0}
	}
}
