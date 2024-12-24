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
	return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
}

type Gate struct {
	operator string
	input1   string
	input2   string
	output   string
}

func formatData(rows []string) (map[string]int, []Gate) {
	wireValuesMap := make(map[string]int)
	gates := []Gate{}

	for _, row := range rows {
		if strings.Contains(row, ":") {
			parts := strings.Split(row, ": ")
			wire := parts[0]
			value, _ := strconv.Atoi(parts[1])
			wireValuesMap[wire] = value

		} else if strings.Contains(row, " -> ") {
			parts := strings.Split(row, " -> ")
			output := parts[1]
			gateParts := strings.Fields(parts[0])

			if len(gateParts) == 3 {
				gates = append(gates, Gate{
					input1:   gateParts[0],
					operator: gateParts[1],
					input2:   gateParts[2],
					output:   output,
				})
			}
		}
	}

	return wireValuesMap, gates
}

func (g Gate) evaluate(in1, in2 int) int {
	switch g.operator {
	case "AND":
		return in1 & in2
	case "OR":
		return in1 | in2
	case "XOR":
		return in1 ^ in2
	default:
		panic("Woot? " + g.operator)
	}
}

func part1(wireValues map[string]int, inputs []Gate) {
	gates := inputs

	for len(gates) > 0 {
		newGates := []Gate{}

		for _, gate := range gates {
			value1, value1Exists := wireValues[gate.input1]
			value2, value2Exists := wireValues[gate.input2]

			// if both values -> store gate's result
			if value1Exists && value2Exists {
				wireValues[gate.output] = gate.evaluate(value1, value2)
			} else {
				newGates = append(newGates, gate)
			}
		}
		gates = newGates
	}

	resultBinary := ""
	for i := 0; ; i++ {
		wireKey := fmt.Sprintf("z%02d", i) // "z00", "z01", "z02", ...

		value, exists := wireValues[wireKey]
		if !exists {
			break
		}
		resultBinary = strconv.Itoa(value) + resultBinary
	}

	resultDecimal, _ := strconv.ParseInt(resultBinary, 2, 64)
	fmt.Println("Part 1:", resultDecimal)
}

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	wireValues, gates := formatData(data)
	part1(wireValues, gates)
	// part2(wireValues, gates)
}
