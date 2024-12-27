package main

import (
	"fmt"
	"os"
	"path/filepath"
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

type System struct {
	registers map[string]int
	program   []int
}

func formatData(rows []string) System {
	system := System{
		registers: make(map[string]int),
		program:   []int{},
	}

	for _, row := range rows {
		if strings.Contains(row, "Register") {
			parts := strings.Split(row, ": ")
			registerName := parts[0]
			registerValue, _ := strconv.Atoi(parts[1])
			system.registers[registerName] = registerValue
		}

		if strings.Contains(row, "Program") {
			parts := strings.Split(row, ": ")
			program := strings.Split(parts[1], ",")

			system.program = make([]int, len(program))
			for i, v := range program {
				system.program[i], _ = strconv.Atoi(v)
			}

		}
	}

	return system
}

func evaluate(system System, opcode int, operand int) ([]string, bool) {
	output := []string{}

	getComboValue := func(val int) int {
		if val <= 3 {
			return val
		}
		switch val {
		case 4:
			return system.registers["Register A"]
		case 5:
			return system.registers["Register B"]
		case 6:
			return system.registers["Register C"]
		default:
			return 0
		}
	}

	switch opcode {
	case 0: // adv - divide A by 2^combo_operand, store in A
		divisor := 1 << getComboValue(operand)
		system.registers["Register A"] = system.registers["Register A"] / divisor

	case 1: // bxl - XOR B with literal operand
		system.registers["Register B"] ^= operand

	case 2: // bst - write combo_operand mod 8 to B
		system.registers["Register B"] = getComboValue(operand) % 8

	case 3: // jnz - jump if A != 0
		if system.registers["Register A"] != 0 {
			return output, true
		}

	case 4: // bxc - XOR B with C
		system.registers["Register B"] ^= system.registers["Register C"]

	case 5: // out - output combo_operand %8
		outVal := getComboValue(operand) % 8
		output = append(output, strconv.Itoa(outVal))

	case 6: // bdv - like adv but store in B
		divisor := 1 << getComboValue(operand)
		system.registers["Register B"] = system.registers["Register A"] / divisor

	case 7: // cdv - like adv but store in C
		divisor := 1 << getComboValue(operand)
		system.registers["Register C"] = system.registers["Register A"] / divisor
	}

	return output, false
}

func execute(inputs System) []string {
	pointer := 0
	outputs := []string{}

	for pointer < len(inputs.program)-1 {
		opcode := inputs.program[pointer]
		operand := inputs.program[pointer+1]

		output, jumped := evaluate(inputs, opcode, operand)
		outputs = append(outputs, output...)

		if jumped {
			pointer = operand
		} else {
			pointer += 2
		}
	}
	return outputs
}

func part1(inputs System) string {
	outputs := execute(inputs)

	result := strings.Join(outputs, ",")
	fmt.Println("Part 1:", result)
	return result
}

func findMinRegisterValue(sys *System, programIndex int, result int) int {
	if programIndex < 0 {
		return result
	}

	for d := 0; d <= 7; d++ {
		testRegisterValue := (result << 3) | d
		i := 0
		var target int

		for i < len(sys.program) {
			var output int
			if sys.program[i+1] <= 3 {
				output = sys.program[i+1]
			} else if sys.program[i+1] == 4 {
				output = testRegisterValue
			} else if sys.program[i+1] == 5 {
				output = sys.registers["B"]
			} else if sys.program[i+1] == 6 {
				output = sys.registers["C"]
			}

			switch sys.program[i] {
			case 0: // adv
				testRegisterValue >>= output
			case 1: // bxl
				sys.registers["B"] ^= sys.program[i+1]
			case 2: // bst
				sys.registers["B"] = output & 7
			case 3: // jnz
				if testRegisterValue != 0 {
					i = sys.program[i+1] - 2
				}
			case 4: // bxc
				sys.registers["B"] ^= sys.registers["C"]
			case 5: // out
				target = output & 7
				goto breakLoop
			case 6: // bdv
				sys.registers["B"] = testRegisterValue >> output
			case 7: // cdv
				sys.registers["C"] = testRegisterValue >> output
			}
			i += 2
		}

	breakLoop:
		if target == sys.program[programIndex] {
			result := findMinRegisterValue(sys, programIndex-1, result<<3|d)
			if result >= 0 {
				return result
			}
		}
	}
	return -1
}

func part2(inputs System) int {
	programLength := len(inputs.program) - 1
	result := findMinRegisterValue(&inputs, programLength, inputs.program[programLength])
	fmt.Println("Part 2:", result)
	return result
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	formattedData := formatData(data)
	part1(formattedData)

	formattedData = formatData(data)
	part2(formattedData)
}
