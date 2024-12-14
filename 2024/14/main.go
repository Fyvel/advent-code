package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func parseRobot(line string) Robot {
	// extract vectors from line such as "p=0,4 v=3,-3"
	re := regexp.MustCompile(`p=([-\d]+),([-\d]+) v=([-\d]+),([-\d]+)`)
	matches := re.FindStringSubmatch(line)
	px, _ := strconv.Atoi(matches[1])
	py, _ := strconv.Atoi(matches[2])
	vx, _ := strconv.Atoi(matches[3])
	vy, _ := strconv.Atoi(matches[4])

	return Robot{
		position: Vector{px, py},
		move:     Vector{vx, vy},
	}
}

func formatData(data []string) map[Vector][]*Robot {
	robots := make(map[Vector][]*Robot)

	for _, row := range data {
		if row == "" {
			continue
		}
		robot := parseRobot(row)
		robots[robot.position] = append(robots[robot.position], &robot)
	}
	return robots
}

type Vector struct {
	x, y int
}
type Robot struct {
	position Vector
	move     Vector
}

func renderGrid(robots map[Vector][]*Robot, width int, height int, hideMid bool) [][]string {
	grid := make([][]string, height)

	for r := range grid {
		grid[r] = make([]string, width)

		for c := range grid[r] {
			if hideMid && c == width/2 || hideMid && r == height/2 {
				grid[r][c] = " "
				continue
			}

			v := Vector{c, r}
			if _, ok := robots[v]; ok {
				grid[r][c] = strconv.Itoa(len(robots[v]))
			} else {
				grid[r][c] = "."
			}
		}
		fmt.Println(grid[r])
	}
	fmt.Println()
	return grid
}

func calculateNextPositions(robotsMap map[Vector][]*Robot, nbOfSeconds int, width int, height int) map[Vector][]*Robot {
	nextRobotsMap := make(map[Vector][]*Robot)
	for _, robots := range robotsMap {
		for _, robot := range robots {
			newPosition := Vector{
				x: ((robot.position.x+(robot.move.x*nbOfSeconds))%width + width) % width,
				y: ((robot.position.y+(robot.move.y*nbOfSeconds))%height + height) % height,
			}
			nextRobotsMap[newPosition] = append(nextRobotsMap[newPosition], robot)
		}
	}
	return nextRobotsMap
}

func getQuadrants(nextRobotsMap map[Vector][]*Robot, midCol int, midRow int) map[string][]*Robot {
	quadrants := make(map[string][]*Robot)

	for v, robots := range nextRobotsMap {
		if v.x < midCol && v.y < midRow {
			quadrants["topLeft"] = append(quadrants["topLeft"], robots...)
		}
		if v.x > midCol && v.y < midRow {
			quadrants["topRight"] = append(quadrants["topRight"], robots...)
		}
		if v.x < midCol && v.y > midRow {
			quadrants["bottomLeft"] = append(quadrants["bottomLeft"], robots...)
		}
		if v.x > midCol && v.y > midRow {
			quadrants["bottomRight"] = append(quadrants["bottomRight"], robots...)
		}
	}
	return quadrants
}

func part1(robotsMap map[Vector][]*Robot) int {
	width := 101
	height := 103
	nbOfSeconds := 100

	nextRobotsMap := calculateNextPositions(robotsMap, nbOfSeconds, width, height)

	midCol := width / 2
	midRow := height / 2

	renderGrid(nextRobotsMap, width, height, true)

	quadrants := getQuadrants(nextRobotsMap, midCol, midRow)

	result := 0
	for _, robots := range quadrants {
		if result == 0 {
			result += len(robots)
			continue
		}
		result *= len(robots)
	}

	fmt.Println("part 1:", result)
	return result
}

func part2(robotsMap map[Vector][]*Robot) {
	if err := os.MkdirAll("dump", 0755); err != nil {
		log.Fatal(err)
	}

	width := 101
	height := 103

	for i := 8000; i <= 60*60*8; i++ {
		nbOfSeconds := i + 1
		nextRobotsMap := calculateNextPositions(robotsMap, nbOfSeconds, width, height)

		quadrants := getQuadrants(nextRobotsMap, width/2, height/2)

		robotCountPerQuadrant := []int{
			len(quadrants["topLeft"]),
			len(quadrants["topRight"]),
			len(quadrants["bottomLeft"]),
			len(quadrants["bottomRight"]),
		}

		renderGrid(nextRobotsMap, width, height, false)
		fmt.Println("Seconds:", nbOfSeconds)
		time.Sleep(150 * time.Millisecond)

		maxRobots := 0
		for _, count := range robotCountPerQuadrant {
			if count > maxRobots {
				maxRobots = count
			}
		}
		if maxRobots > len(robotsMap)/2 {
			renderGrid(nextRobotsMap, width, height, false)
			timeFormatted := fmt.Sprintf("%02dh:%02dm:%02ds", nbOfSeconds/3600, (nbOfSeconds/60)%60, nbOfSeconds%60)
			fmt.Printf("Part 2: %d (%s)", nbOfSeconds, timeFormatted)
			break
		}
	}
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
