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

func main() {
	data, err := readData()
	if err != nil {
		fmt.Println("Get rekt:", err)
		return
	}

	// fmt.Println(data)
	part1(data)
	part2(data)
}

func part1(data []string) {
	count := 0
	w := NewWheel()

	for _, row := range data {
		direction, distance := parseData(row)

		switch direction {
		case 'L':
			w.decrement(distance % w.max)
		case 'R':
			w.increment(distance % w.max)
		}

		if w.pos == w.min {
			count++
		}
	}

	fmt.Println("Password: ", count)
}

func parseData(row string) (byte, int) {
	direction := row[0]
	distance, _ := strconv.Atoi(row[1:])
	return direction, distance
}

func part2(data []string) {
	w := NewWheel()

	for _, row := range data {
		direction, distance := parseData(row)

		switch direction {
		case 'L':
			w.decrement(distance)
		case 'R':
			w.increment(distance)
		}
	}

	fmt.Println("Password: ", w.carry)
}

func NewWheel() *Wheel {
	return &Wheel{
		pos: 50,
		min: 0,
		max: 100,
	}
}

type Wheel struct {
	pos   int
	min   int
	max   int
	carry int
}

func (w *Wheel) decrement(dist int) {
	for range dist {
		w.pos--
		if w.pos < w.min {
			w.pos = w.max - 1
		}
		if w.pos == 0 {
			w.carry++
		}
	}
}

func (w *Wheel) increment(dist int) {
	for range dist {
		w.pos++
		if w.pos >= w.max {
			w.pos = 0
		}
		if w.pos == 0 {
			w.carry++
		}
	}
}
