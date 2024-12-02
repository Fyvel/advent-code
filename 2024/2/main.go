import (
	"fmt"
	"math"
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

func formatData(reports []string) [][]int {
	levels := make([][]int, len(reports))
	for i, report := range reports {
		nums := strings.Fields(report)
		level := make([]int, len(nums))
		for j, num := range nums {
			var n int
			fmt.Sscanf(num, "%d", &n)
			level[j] = n
		}
		levels[i] = level
	}
	return levels
}

func part1(levels [][]int) int {
	safeLevelSeen := 0
	for _, level := range levels {
		if checkLevel(level) {
			safeLevelSeen++
		}
	}
	fmt.Printf("safeLevelSeen: %d\n", safeLevelSeen)
	return safeLevelSeen
}

func part2(levels [][]int) int {
	memo := make(map[string]bool)
	safeLevelSeen := 0

	for _, level := range levels {
		levelKey := fmt.Sprint(level)
		if val, exists := memo[levelKey]; exists {
			if val {
				safeLevelSeen++
			}
			continue
		}

		isSafeLevel := checkLevel(level)
		memo[levelKey] = isSafeLevel

		if !isSafeLevel {
			isSafeLevel = retry(level, memo)
		}

		if isSafeLevel {
			safeLevelSeen++
		}
	}
	fmt.Printf("safeLevelSeen: %d\n", safeLevelSeen)
	return safeLevelSeen
}

func checkLevel(level []int) bool {
	isSafeLevel := true
	direction := ""

	for j := 0; j < len(level)-1; j++ {
		var currentDirection string
		if level[j] >= level[j+1] {
			currentDirection = "dec"
		} else {
			currentDirection = "inc"
		}

		if direction != "" && direction != currentDirection {
			isSafeLevel = false
			break
		}
		direction = currentDirection

		gap := math.Abs(float64(level[j] - level[j+1]))
		inRange := gap >= 1 && gap <= 3
		if !inRange {
			isSafeLevel = false
			break
		}
	}
	return isSafeLevel
}

func retry(level []int, memo map[string]bool) bool {
	var isSafeLevel bool
	for j := 0; j < len(level); j++ {
		newLevel := append(level[:j], level[j+1:]...)
		newLevelKey := fmt.Sprint(newLevel)

		if val, exists := memo[newLevelKey]; exists {
			isSafeLevel = val
			if isSafeLevel {
				break
			}
		}

		isSafeLevel = checkLevel(newLevel)
		memo[newLevelKey] = isSafeLevel

		if isSafeLevel {
			break
		}
	}
	return isSafeLevel
}

func main() {
	reports, err := readData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formattedData := formatData(reports)
	part1(formattedData)
	part2(formattedData)
}
