package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

type DoorCode struct {
	code         string
	num          int
	instructions string
}

func formatData(rows []string) []DoorCode {
	doorCodes := make([]DoorCode, len(rows))
	regex := regexp.MustCompile(`\d+`)
	for r, row := range rows {
		// extract number from row removing leading zeros
		num, _ := strconv.Atoi(regex.FindString(row))
		doorCodes[r] = DoorCode{code: row, num: num, instructions: ""}
	}
	return doorCodes
}

func (keypad *Keypad) generateInstructions(code string) string {
	instructions := ""
	currentPosition := keypad.position

	for _, char := range code {
		// add path to instruction
		instructions += keypad.pathFromTo[currentPosition][string(char)]
		// update position
		currentPosition = string(char)
	}
	return instructions
}

type Position [2]int

type KeypadPathMap map[string]map[string]string

type Keypad struct {
	keypad     [][]string
	position   string
	directions [4]Position
	pathFromTo KeypadPathMap
}

func getNumericKeypad() *Keypad {
	return &Keypad{
		keypad: [][]string{
			{"7", "8", "9"},
			{"4", "5", "6"},
			{"1", "2", "3"},
			{"", "0", "A"},
		},
		position: "A",
		directions: [4]Position{
			{-1, 0}, // ^
			{0, 1},  // >
			{1, 0},  // v
			{0, -1}, // <
		},
		pathFromTo: KeypadPathMap{
			"A": {
				"0": "<A",
				"7": "^^^<<A",
				"8": "<^^^A",
				"9": "^^^A",
				"4": "^^<<A",
				"5": "<^^A",
				"6": "^^A",
				"1": "^<<A",
				"2": "<^A",
				"3": "^A",
				"A": "A",
			},
			"0": {
				"0": "A",
				"7": "^^^<A",
				"8": "^^^A",
				"9": "^^^>A",
				"4": "^<^A",
				"5": "^^A",
				"6": "^^>A",
				"1": "^<A",
				"2": "^A",
				"3": "^>A",
				"A": ">A",
			},
			"1": {
				"0": ">vA",
				"7": "^^A",
				"8": "^^>A",
				"9": "^^>>A",
				"4": "^A",
				"5": "^>A",
				"6": "^>>A",
				"1": "A",
				"2": ">A",
				"3": ">>A",
				"A": ">>vA",
			},
			"2": {
				"0": "vA",
				"7": "<^^A",
				"8": "^^A",
				"9": "^^>A",
				"4": "<^A",
				"5": "^A",
				"6": "^>A",
				"1": "<A",
				"2": "A",
				"3": ">A",
				"A": "v>A",
			},
			"3": {
				"0": "<vA",
				"7": "<<^^A",
				"8": "<^^A",
				"9": "^^A",
				"4": "<<^A",
				"5": "<^A",
				"6": "^A",
				"1": "<<A",
				"2": "<A",
				"3": "A",
				"A": "vA",
			},
			"4": {
				"0": ">vvA",
				"7": "^A",
				"8": "^>A",
				"9": "^>>A",
				"4": "A",
				"5": ">A",
				"6": ">>A",
				"1": "vA",
				"2": "v>A",
				"3": "v>>A",
				"A": ">>vvA",
			},
			"5": {
				"0": "vvA",
				"7": "<^A",
				"8": "^A",
				"9": "^>A",
				"4": "<A",
				"5": "A",
				"6": ">A",
				"1": "<vA",
				"2": "vA",
				"3": "v>A",
				"A": "vv>A",
			},
			"6": {
				"0": "<vvA",
				"7": "<<^A",
				"8": "<^A",
				"9": "^A",
				"4": "<<A",
				"5": "<A",
				"6": "A",
				"1": "<<vA",
				"2": "<vA",
				"3": "vA",
				"A": "vvA",
			},
			"7": {
				"0": ">vvvA",
				"7": "A",
				"8": ">A",
				"9": ">>A",
				"4": "vA",
				"5": "v>A",
				"6": "v>>A",
				"1": "vvA",
				"2": "vv>A",
				"3": "vv>>A",
				"A": ">>vvvA",
			},
			"8": {
				"0": "vvvA",
				"7": "<A",
				"8": "A",
				"9": ">A",
				"4": "<vA",
				"5": "^A",
				"6": "v>A",
				"1": "<vvA",
				"2": "vvA",
				"3": "vv>A",
				"A": "vvv>A",
			},
			"9": {
				"0": "<vvvA",
				"7": "<<A",
				"8": "<A",
				"9": "A",
				"4": "<<vA",
				"5": "<vA",
				"6": "vA",
				"1": "<<vvA",
				"2": "<vvA",
				"3": "vvA",
				"A": "vvvA",
			},
		},
	}
}

func getDirectionalKeypad() *Keypad {
	return &Keypad{
		keypad: [][]string{
			{"", "^", "A"},
			{"<", "v", ">"},
		},
		position: "A",
		directions: [4]Position{
			{1, 0},  // v
			{0, 1},  // >
			{-1, 0}, // ^
			{0, -1}, // <
		},
		pathFromTo: KeypadPathMap{
			"A": {
				"^": "<A",
				">": "vA",
				"v": "<vA",
				"<": "v<<A",
				"A": "A",
			},
			"<": {
				"^": ">^A",
				"v": ">A",
				">": ">>A",
				"A": ">>^A",
				"<": "A",
			},
			"v": {
				"^": "^A",
				"<": "<A",
				">": ">A",
				"A": "^>A",
				"v": "A",
			},
			">": {
				"v": "<A",
				"<": "<<A",
				"A": "^A",
				"^": "<^A",
				">": "A",
			},
			"^": {
				"v": "vA",
				"<": "v<A",
				"A": ">A",
				">": "v>A",
				"^": "A",
			},
		},
	}
}

func part1(doorCodes []DoorCode) int {
	sum := 0

	for _, doorCode := range doorCodes {
		doorInstruction := getNumericKeypad().generateInstructions(doorCode.code)
		robotInstruction := getDirectionalKeypad().generateInstructions(doorInstruction)
		shortestInstruction := getDirectionalKeypad().generateInstructions(robotInstruction)

		fmt.Println(doorCode.code, ":", len(shortestInstruction), "x", doorCode.num)
		sum += len(shortestInstruction) * doorCode.num
	}

	fmt.Println("Part 1:", sum)
	return sum
}

type MemoKey struct {
	path  string
	level int
}

// store the length directly since concatenated string are too long
var memo = make(map[MemoKey]int)

func (k *Keypad) countInstructionsLength(instructionsPath string, depth int) int {
	key := MemoKey{path: instructionsPath, level: depth}
	if value, exists := memo[key]; exists {
		return value
	}

	length := 0
	if depth == 0 {
		length = len(instructionsPath)
		memo[key] = length
		return length
	}

	currentPosition := k.position
	for _, char := range instructionsPath {
		if currentPosition == string(char) {
			length += len(k.pathFromTo[currentPosition][string(char)])
			currentPosition = string(char)
			continue
		}

		path := k.pathFromTo[currentPosition][string(char)]
		length += k.countInstructionsLength(path, depth-1)
		currentPosition = string(char)
	}

	memo[key] = length
	return length
}

func part2(doorCodes []DoorCode) int {
	sum := 0
	numberOfRobots := 25

	for _, doorCode := range doorCodes {
		doorInstruction := getNumericKeypad().generateInstructions(doorCode.code)
		instruction := getDirectionalKeypad().countInstructionsLength(doorInstruction, numberOfRobots)
		fmt.Println(doorCode.code, ":", instruction, "x", doorCode.num)
		sum += instruction * doorCode.num
	}

	fmt.Println("Part 2:", sum)
	return sum
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
