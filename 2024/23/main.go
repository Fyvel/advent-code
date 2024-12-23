package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

type Graph map[string]map[string]bool

func formatData(rows []string) Graph {
	graph := make(Graph)

	for _, row := range rows {
		nodes := strings.Split(row, "-")
		if len(nodes) != 2 {
			continue
		}

		if _, exists := graph[nodes[0]]; !exists {
			graph[nodes[0]] = make(map[string]bool)
		}
		if _, exists := graph[nodes[1]]; !exists {
			graph[nodes[1]] = make(map[string]bool)
		}

		graph[nodes[0]][nodes[1]] = true
		graph[nodes[1]][nodes[0]] = true
	}

	return graph
}

func (graph Graph) findTriangles() [][]string {
	triangles := make([][]string, 0)
	seenTriangles := make(map[string]bool)

	for node1 := range graph {
		for node2 := range graph[node1] {
			if node2 == node1 {
				continue
			}
			for node3 := range graph[node2] {
				if _, exists := graph[node3][node1]; exists {
					triangle := []string{node1, node2, node3}
					sort.Strings(triangle)
					key := strings.Join(triangle, ",")

					if _, exists := seenTriangles[key]; !exists {
						triangles = append(triangles, triangle)
						seenTriangles[key] = true
					}
				}
			}
		}
	}
	return triangles
}

func part1(graph Graph) int {
	triangles := graph.findTriangles()
	seenTriangles := make(map[string]bool)
	count := 0

	for _, triangle := range triangles {
		key := strings.Join(triangle, "-")

		for _, node := range triangle {
			if strings.HasPrefix(node, "t") {
				seenTriangles[key] = true
				count++
				break
			}
		}
	}

	fmt.Printf("Part 1: %d triangles starting with 't'\n", count)
	return count
}

func (graph Graph) findLongestConnection(triangles [][]string) []string {
	largestGroup := []string{}

	isConnectedToAll := func(node string, group []string) bool {
		for _, member := range group {
			if !graph[node][member] {
				return false
			}
		}
		return true
	}

	for _, triangle := range triangles {
		currentGroup := make([]string, len(triangle))
		copy(currentGroup, triangle)

		// expand group
		for node := range graph {
			if contains(currentGroup, node) {
				continue
			}

			// check inter-connection
			if isConnectedToAll(node, currentGroup) {
				currentGroup = append(currentGroup, node)
			}
		}

		// update largest group
		if len(currentGroup) > len(largestGroup) {
			largestGroup = currentGroup
		}
	}

	return largestGroup
}

func contains(slice []string, str string) bool {
	for _, val := range slice {
		if val == str {
			return true
		}
	}
	return false
}

func part2(graph Graph) string {
	triangles := graph.findTriangles()

	longestConnections := graph.findLongestConnection(triangles)
	sort.Strings(longestConnections)

	longestConnection := strings.Join(longestConnections, ",")
	fmt.Printf("Part 2: %s\n", longestConnection)
	return longestConnection
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
