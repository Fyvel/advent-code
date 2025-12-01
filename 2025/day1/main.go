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
}

func part1(data []string) {
	count := 0
	w := Wheel{
		pos: 50,
		min: 0,
		max: 100,
	}

	for _, row := range data {
		direction := row[0]
		distance, _ := strconv.Atoi(row[1:])

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

type Wheel struct {
	pos int
	min int
	max int
}

func (w *Wheel) decrement(dist int) {
	if (w.pos - dist) < w.min {
		w.pos = w.pos + w.max - dist
		return
	}
	w.pos -= dist
}

func (w *Wheel) increment(dist int) {
	if (w.pos + dist) >= w.max {
		w.pos = w.pos - w.max + dist
		return
	}
	w.pos += dist

}
