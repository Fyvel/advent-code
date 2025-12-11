package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	Reset       = "\033[0m"
	Yellow      = "\033[33m"
	Green       = "\033[32m"
	Black       = "\033[90m"
	BgGreen     = "\033[42m"
	BgWhite     = "\033[47m"
	BgOrange    = "\033[48;5;208m"
	ClearScreen = "\033[2J"
	MoveCursor  = "\033[H"
)

type Visualiser struct {
	spinnerChars []string
	spinnerIdx   int
	delay        time.Duration
	disabled     bool
	mu           sync.Mutex
	lineMap      map[int]int // maps lights index to terminal line
	totalLines   int
}

func NewVisualiser(delay time.Duration, disabled bool) *Visualiser {
	return &Visualiser{
		spinnerChars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		spinnerIdx:   0,
		delay:        delay,
		disabled:     disabled,
		lineMap:      make(map[int]int),
		totalLines:   0,
	}
}

func renderLight(lights string, isOn bool) string {
	if isOn {
		return fmt.Sprintf("%s 🟡 %s[%s]%s", BgGreen, Black, lights, Reset)
	}
	return fmt.Sprintf("%s ⚫️ %s[%s]%s ", BgWhite, Black, lights, Reset)
}

func renderButtons(buttons [][]int, active int) string {
	var sb strings.Builder
	for i, button := range buttons {
		schema := make([]rune, len(buttons))
		for j := range schema {
			schema[j] = '.'
		}
		for _, idx := range button {
			if idx >= 0 && idx < len(schema) {
				schema[idx] = '#'
			}
		}
		if i == active {
			sb.WriteString(fmt.Sprintf("%s(%s)%s ", BgOrange, string(schema), Reset))
		} else {
			sb.WriteString(fmt.Sprintf("(%s) ", string(schema)))
		}
	}
	return sb.String()
}

func (v *Visualiser) RegisterMachine(machineIdx int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.lineMap[machineIdx] = v.totalLines
	v.totalLines++
}

func (v *Visualiser) Render(light string, isOn bool, buttons [][]int, active int) {
	if v.disabled {
		return
	}
	fmt.Print(ClearScreen)
	fmt.Print(MoveCursor)
	fmt.Printf("%s %s\n", renderLight(light, isOn), renderButtons(buttons, active))
	time.Sleep(v.delay)
}

func (v *Visualiser) moveTo(line int) {
	if line > 0 {
		fmt.Printf("\033[%dA", line) // Move cursor up
	}
}

func (v *Visualiser) moveDown(line int) {
	if line > 0 {
		fmt.Printf("\033[%dB", line) // Move cursor down
	}
}

func (v *Visualiser) Update(machineIdx int, light string, isOn bool, buttons [][]int, active int) {
	if v.disabled {
		return
	}
	v.mu.Lock()
	spinner := v.spinnerChars[v.spinnerIdx%len(v.spinnerChars)]
	v.spinnerIdx++
	line := v.lineMap[machineIdx]
	currentLine := v.totalLines - 1
	v.mu.Unlock()

	v.mu.Lock()
	v.moveTo(currentLine - line)
	fmt.Printf("\r%s %s %s\033[K", spinner, renderLight(light, isOn), renderButtons(buttons, active))
	v.moveDown(currentLine - line)
	v.mu.Unlock()

	time.Sleep(v.delay)
}

func (v *Visualiser) Complete(machineIdx int, light string, buttons [][]int) {
	if v.disabled {
		return
	}
	v.mu.Lock()
	line := v.lineMap[machineIdx]
	currentLine := v.totalLines - 1
	v.mu.Unlock()

	v.mu.Lock()
	v.moveTo(currentLine - line)
	fmt.Printf("\r✓ %s %s\033[K", renderLight(light, true), renderButtons(buttons, -1))
	v.moveDown(currentLine - line)
	v.mu.Unlock()

	time.Sleep(v.delay)
}
