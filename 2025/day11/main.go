package main

import (
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

func formatData(rows []string) map[string][]string {
	hash := make(map[string][]string)
	for _, row := range rows {
		parts := strings.SplitN(row, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		values := strings.Split(strings.TrimSpace(parts[1]), " ")
		for val := range values {
			values[val] = strings.TrimSpace(values[val])
		}
		hash[key] = values
	}
	return hash
}

type TreeNode struct {
	Value    string
	Children []*TreeNode
}

func part1(data map[string][]string) {
	// build Tree
	root := &TreeNode{Value: "you"}
	nodes := map[string]*TreeNode{
		"you": root,
	}

	for parent, children := range data {
		parentNode, exists := nodes[parent]
		if !exists {
			parentNode = &TreeNode{Value: parent}
			nodes[parent] = parentNode
		}
		for _, child := range children {
			childNode, exists := nodes[child]
			if !exists {
				childNode = &TreeNode{Value: child}
				nodes[child] = childNode
			}
			parentNode.Children = append(parentNode.Children, childNode)
		}
	}

	// traverse tree
	memo := make(map[*TreeNode][][]string)
	visited := make(map[*TreeNode]bool)
	paths := traverseTreeDFS(root, []string{}, visited, memo)
	for _, path := range paths {
		fmt.Println(strings.Join(path, " -> "))
	}

	fmt.Println("Part 1:", len(paths))
}

func traverseTreeDFS(node *TreeNode, currentPath []string, visited map[*TreeNode]bool, memo map[*TreeNode][][]string) [][]string {
	if visited[node] {
		return nil
	}

	visited[node] = true
	defer delete(visited, node)
	if cached, exists := memo[node]; exists {
		result := make([][]string, len(cached))
		for i, cachedPath := range cached {
			result[i] = make([]string, len(currentPath)+len(cachedPath))
			copy(result[i], currentPath)
			copy(result[i][len(currentPath):], cachedPath)
		}
		return result
	}

	currentPath = append(currentPath, node.Value)

	// base case
	if len(node.Children) == 0 && node.Value == "out" {
		pathCopy := make([]string, len(currentPath))
		copy(pathCopy, currentPath)
		return [][]string{pathCopy}
	}

	// recursion
	var paths [][]string
	for _, child := range node.Children {
		childPaths := traverseTreeDFS(child, currentPath, visited, memo)
		paths = append(paths, childPaths...)
	}

	// store
	if len(paths) > 0 {
		nodePaths := make([][]string, len(paths))
		for i, path := range paths {
			//  part after current node
			nodePaths[i] = path[len(currentPath):]
		}
		memo[node] = nodePaths
	}

	return paths
}

func part2(data map[string][]string) {
	root := &TreeNode{Value: "svr"}
	nodes := map[string]*TreeNode{
		"svr": root,
	}

	for parent, children := range data {
		parentNode, exists := nodes[parent]
		if !exists {
			parentNode = &TreeNode{Value: parent}
			nodes[parent] = parentNode
		}
		for _, child := range children {
			childNode, exists := nodes[child]
			if !exists {
				childNode = &TreeNode{Value: child}
				nodes[child] = childNode
			}
			parentNode.Children = append(parentNode.Children, childNode)
		}
	}

	visited := make(map[*TreeNode]bool)

	memo := make(map[MemoKey]int)
	count := countPathsDFS(root, visited, false, false, memo)

	fmt.Println("Part 2:", count)
}

type MemoKey struct {
	node   *TreeNode
	hasDac bool
	hasFft bool
}

func countPathsDFS(node *TreeNode, visited map[*TreeNode]bool, hasDac bool, hasFft bool, memo map[MemoKey]int) int {
	// base
	if visited[node] {
		return 0
	}

	key := MemoKey{
		node:   node,
		hasDac: hasDac,
		hasFft: hasFft,
	}

	// check memo
	if count, exists := memo[key]; exists {
		return count
	}

	visited[node] = true
	// cleanup
	defer delete(visited, node)

	switch node.Value {
	case "dac":
		hasDac = true
	case "fft":
		hasFft = true
	case "out":
		if hasDac && hasFft {
			pathStr := []string{}
			for n := range visited {
				pathStr = append(pathStr, n.Value)
			}
			pathStr = append(pathStr, "out")
			fmt.Println(strings.Join(pathStr, " -> "))
			return 1
		}
		return 0
	}

	// recursion
	count := 0
	for _, child := range node.Children {
		count += countPathsDFS(child, visited, hasDac, hasFft, memo)
	}

	// store
	memo[key] = count
	return count
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
