package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	colourReset  = "\033[0m"
	colourYellow = "\033[33m"
	colourGreen  = "\033[32m"
	colourBlack  = "\033[90m"
)

type Visualiser struct {
	spinnerChars []string
	spinnerIdx   int
	delay        time.Duration
	disabled     bool
	mu           sync.Mutex
	lineMap      map[int]int // maps bank index to terminal line
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

func (v *Visualiser) RegisterBank(bankIdx int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.lineMap[bankIdx] = v.totalLines
	v.totalLines++
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

func (v *Visualiser) UpdateSearching(bankIdx int, bank string, startIdx, currentIdx, maxDigitIdx int, joltage string) {
	if v.disabled {
		return
	}
	v.mu.Lock()
	spinner := v.spinnerChars[v.spinnerIdx%len(v.spinnerChars)]
	v.spinnerIdx++
	line := v.lineMap[bankIdx]
	currentLine := v.totalLines - 1
	v.mu.Unlock()

	colourisedBank := ""
	for i, char := range bank {
		if i < startIdx {
			continue
		} else if i == currentIdx {
			colourisedBank += colourYellow + string(char) + colourReset
		} else if i == maxDigitIdx && maxDigitIdx >= startIdx {
			colourisedBank += colourGreen + string(char) + colourReset
		} else {
			colourisedBank += colourBlack + string(char) + colourReset
		}
	}

	// joltage in green
	colourisedJoltage := colourGreen + joltage + colourReset

	v.mu.Lock()
	v.moveTo(currentLine - line)
	fmt.Printf("\r%s 🪫%s <- %s\033[K", spinner, colourisedJoltage, colourisedBank)
	v.moveDown(currentLine - line)
	v.mu.Unlock()

	time.Sleep(v.delay)
}

func (v *Visualiser) Update(bank, joltage string, maxLength int, usedIndices map[int]bool) {
	if v.disabled {
		return
	}
	spinner := v.spinnerChars[v.spinnerIdx%len(v.spinnerChars)]
	v.spinnerIdx++

	colourisedBank := ""
	for i, char := range bank {
		if usedIndices[i] {
			// current digit in yellow
			colourisedBank += colourYellow + string(char) + colourReset
		} else {
			// discarded digits in black
			colourisedBank += colourBlack + string(char) + colourReset
		}
	}

	// joltage in green
	colourisedJoltage := colourGreen + joltage + colourReset
	padding := strings.Repeat(" ", maxLength-len(joltage))
	fmt.Printf("\r%s %s -> %s%s", spinner, colourisedBank, colourisedJoltage, padding)
	time.Sleep(v.delay)
}

func (v *Visualiser) Complete(bankIdx int, bank, joltage string, usedIndices map[int]bool) {
	if v.disabled {
		return
	}
	v.mu.Lock()
	line := v.lineMap[bankIdx]
	currentLine := v.totalLines - 1
	v.mu.Unlock()

	// yellow for used, black for discarded
	colourisedBank := ""
	for i, char := range bank {
		if usedIndices[i] {
			colourisedBank += colourYellow + string(char) + colourReset
		} else {
			colourisedBank += colourBlack + string(char) + colourReset
		}
	}

	// joltage in green
	colourisedJoltage := colourGreen + joltage + colourReset

	v.mu.Lock()
	v.moveTo(currentLine - line)
	fmt.Printf("\r✓ 🔋%s <- %s\033[K", colourisedJoltage, colourisedBank)
	v.moveDown(currentLine - line)
	v.mu.Unlock()
}
