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

func (keypad *Keypad) generateInstructions(code string) []string {
	instructions := []string{""}
	currentPosition := keypad.position

	for _, char := range code {
		// get all possible paths
		paths := keypad.pathFromTo[currentPosition][string(char)]

		// add each path at the end of each instruction
		newInstructions := []string{}
		for _, instruction := range instructions {
			for _, path := range paths {
				newInstructions = append(newInstructions, instruction+path)
			}
		}
		instructions = newInstructions
		currentPosition = string(char)
	}
	return instructions
}

type Position [2]int

type Keypad struct {
	keypad     [][]string
	position   string
	directions [4]Position
	pathFromTo map[string]map[string][]string
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
		pathFromTo: map[string]map[string][]string{
			"A": {
				"0": {"<A"},
				"1": {"^<<A"}, // , "<^<A"},
				"2": {"^<A", "<^A"},
				"3": {"^A"},
				"4": {"^^<<A"}, // , "^<^<A", "^<<^A", "<^<^A", "<^^<A"},
				"5": {"^^<A", "<^^A"},
				"6": {"^^A"},
				"7": {"^^^<<A"},         // , "^^<^<A", "^<<^^A", "^<^<^A", "^<^^<A", "<^<^^A", "^^<<^A", "<^^<^A", "<^^^<A"},
				"8": {"^^^<A", "<^^^A"}, // , "^<^^A", "^^<^A"},
				"9": {"^^^A"},
				"A": {"A"},
			},
			"0": {
				"0": {"A"},
				"1": {"^<A"},
				"2": {"^A"},
				"3": {"^>A", ">^A"},
				"4": {"^^<A"}, //  "^<^A"},
				"5": {"^^A"},
				"6": {"^^>A", ">^^A"}, // , "^>^A"},
				"7": {"^^^<A"},        // , "^^<^A", "^<^^A"},
				"8": {"^^^A"},
				"9": {"^^^>A", ">^^^A"}, // , "^>^^A", "^^>^A"},
				"A": {">A"},
			},
			"1": {
				"0": {">vA"},
				"7": {"^^A"},
				"8": {"^^>A", ">^^A"},   // , "^>^A"},
				"9": {"^^>>A", ">>^^A"}, // , ">^>^A", "^>^>A", "^>>^A", ">^^>A"},
				"4": {"^A"},
				"5": {"^>A", ">^A"},
				"6": {"^>>A", ">>^A"}, // , ">^>A"},
				"1": {"A"},
				"2": {">A"},
				"3": {">>A"},
				"A": {">>vA"}, // , ">v>A"},
			},
			"2": {
				"0": {"vA"},
				"7": {"^^<A", "<^^A"}, // , "^<^A"},
				"8": {"^^A"},
				"9": {"^^>A", ">^^A"}, // , "^>^A"},
				"4": {"^<A", "<^A"},
				"5": {"^A"},
				"6": {"^>A", ">^A"},
				"1": {"<A"},
				"2": {"A"},
				"3": {">A"},
				"A": {">vA", "v>A"},
			},
			"3": {
				"0": {"<vA", "v<A"},
				"7": {"^^<<A", "<<^^A"}, // , "<^<^A", "^<^<A", "^<<^A"},
				"8": {"^^<A", "<^^A"},   // , "^<^A"},
				"9": {"^^A"},
				"4": {"^<<A", "<<^A"}, // , "<^<A"},
				"5": {"^<A", "<^A"},
				"6": {"^A"},
				"1": {"<<A"},
				"2": {"<A"},
				"3": {"A"},
				"A": {"vA"},
			},
			"4": {
				"0": {">vvA"}, // , "v>vA"},
				"7": {"^A"},
				"8": {"^>A", ">^A"},
				"9": {">>^A", "^>>A"}, // , ">^>A"},
				"4": {"A"},
				"5": {">A"},
				"6": {">>A"},
				"1": {"vA"},
				"2": {"v>A", ">vA"},
				"3": {">>vA", "v>>A"}, // , ">v>A"},
				"A": {">>vvA"},        // , ">v>vA", "v>>vA", "v>v>A", ">vv>A"},
			},
			"5": {
				"0": {"vvA"},
				"7": {"^<A", "<^A"},
				"8": {"^A"},
				"9": {"^>A", ">^A"},
				"4": {"<A"},
				"5": {"A"},
				"6": {">A"},
				"1": {"v<A", "<vA"},
				"2": {"vA"},
				"3": {"v>A", ">vA"},
				"A": {">vvA", "vv>A"}, // , "v>vA"},
			},
			"6": {
				"0": {"<vvA", "vv<A"}, // , "v<vA"},
				"7": {"<<^A", "^<<A"}, // , "<^<A"},
				"8": {"^<A", "<^A"},
				"9": {"^A"},
				"4": {"<<A"},
				"5": {"<A"},
				"6": {"A"},
				"1": {"<<vA", "v<<A"}, // , "<v<A"},
				"2": {"v<A"},          // , "<vA"},
				"3": {"vA"},
				"A": {"vvA"},
			},
			"7": {
				"0": {">vvvA"}, // , "v>vvA", "vv>vA"},
				"7": {"A"},
				"8": {">A"},
				"9": {">>A"},
				"4": {"vA"},
				"5": {">vA", "v>A"},
				"6": {">>vA", "v>>A"}, // , ">v>A"},
				"1": {"vvA"},
				"2": {">vvA", "vv>A"},   // , "v>vA"},
				"3": {">>vvA", "vv>>A"}, // , ">v>vA", "v>v>A", "v>>vA"},
				"A": {">>vvvA"},         // , ">v>vvA", "vv>>vA", "v>v>vA", ">vv>vA", "v>vv>A", "v>>vvA", "vv>v>A", ">vvv>A"},
			},
			"8": {
				"0": {"vvvA"},
				"7": {"<A"},
				"8": {"A"},
				"9": {">A"},
				"4": {"v<A", "<vA"},
				"5": {"^A"},
				"6": {"v>A", ">vA"},
				"1": {"vv<A", "<vvA"}, // , "v<vA"},
				"2": {"vvA"},
				"3": {"vv>A", ">vvA"},   // , "v>vA"},
				"A": {">vvvA", "vvv>A"}, // , "v>vvA", "vv>vA"},
			},
			"9": {
				"0": {"<vvvA", "vvv<A"}, // , "vv<vA", "v<vvA"},
				"7": {"<<A"},
				"8": {"<A"},
				"9": {"A"},
				"4": {"<<vA", "v<<A"}, // , "<v<A"},
				"5": {"<vA", "v<A"},
				"6": {"vA"},
				"1": {"<<vvA", "vv<<A"}, // , "v<v<A", "<v<vA", "v<<vA", "<vv<A"},
				"2": {"<vvA", "vv<A"},   // , "v<vA"},
				"3": {"vvA"},
				"A": {"vvvA"},
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
		pathFromTo: map[string]map[string][]string{
			"A": {
				"^": {"<A"},
				">": {"vA"},
				"v": {"<vA"},  // , "v<A"},
				"<": {"v<<A"}, // , "<v<A"},
				"A": {"A"},
			},
			"<": {
				"^": {">^A"},
				"v": {">A"},
				">": {">>A"},
				"A": {">>^A"}, // , ">^>A"},
				"<": {"A"},
			},
			"v": {
				"^": {"^A"},
				"<": {"<A"},
				">": {">A"},
				"A": {">^A"}, // , "^>A"},
				"v": {"A"},
			},
			">": {
				"v": {"<A"},
				"<": {"<<A"},
				"A": {"^A"},
				"^": {"^<A"}, // , "<^A"},
				">": {"A"},
			},
			"^": {
				"v": {"vA"},
				"<": {"v<A"},
				"A": {">A"},
				">": {">vA"}, // , "v>A"},
				"^": {"A"},
			},
		},
	}
}

func part1(doorCodes []DoorCode) int {
	sum := 0
	instructions := []string{}

	for _, doorCode := range doorCodes {
		shortestInstruction := ""
		doorInstructions := getNumericKeypad().generateInstructions(doorCode.code)
		for _, doorInstruction := range doorInstructions {
			robotInstructions := getDirectionalKeypad().generateInstructions(doorInstruction)
			for _, robotInstruction := range robotInstructions {
				instructions = getDirectionalKeypad().generateInstructions(robotInstruction)
			}
		}

		// shortest instruction
		for _, instruction := range instructions {
			if shortestInstruction == "" || len(instruction) < len(shortestInstruction) {
				shortestInstruction = instruction
			}
		}

		fmt.Println(doorCode.code, ":", len(shortestInstruction), "x", doorCode.num)
		sum += len(shortestInstruction) * doorCode.num
	}

	fmt.Println("Part 1:", sum)
	return sum
}

func part2(doorCodes []DoorCode) int {
	sum := 0
	// numberOfRobots := 25

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
