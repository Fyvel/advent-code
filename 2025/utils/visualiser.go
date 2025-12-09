package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	Reset        = "\033[0m"
	White        = "\033[97m"
	Black        = "\033[30m"
	Orange       = "\033[38;5;208m"
	Cyan         = "\033[96m"
	BgBlack      = "\033[40m"
	BgOrange     = "\033[48;5;208m"
	BgCyan       = "\033[48;5;51m"
	BgRed        = "\033[41m"
	ClearScreen  = "\033[2J"
	MoveCursor   = "\033[H"
	AltScreenOn  = "\033[?1049h"
	AltScreenOff = "\033[?1049l"
	HideCursor   = "\033[?25l"
	ShowCursor   = "\033[?25h"
)

type CellRenderContext struct {
	Cell           string
	IsActive       bool
	IsInActivePath bool
}

type CellRenderer func(ctx CellRenderContext) string

func RenderGrid(grid [][]string, activeRow, activeCol int, activePath map[string]bool, cellRenderer CellRenderer) {
	if cellRenderer == nil {
		cellRenderer = func(ctx CellRenderContext) string {
			return ctx.Cell
		}
	}

	newLines := make([]string, len(grid))
	for rowIdx, row := range grid {
		var sb strings.Builder
		for colIdx, cell := range row {
			isActive := rowIdx == activeRow && colIdx == activeCol
			key := fmt.Sprintf("%d_%d", rowIdx, colIdx)
			isInActivePath := activePath != nil && activePath[key]

			ctx := CellRenderContext{
				Cell:           cell,
				IsActive:       isActive,
				IsInActivePath: isInActivePath,
			}
			sb.WriteString(cellRenderer(ctx))
		}
		newLines[rowIdx] = sb.String()
	}

	if prevLines == nil || len(prevLines) != len(newLines) {
		var buf strings.Builder
		buf.WriteString(ClearScreen)
		buf.WriteString(MoveCursor)
		for _, line := range newLines {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		fmt.Print(buf.String())
	} else {
		for i, line := range newLines {
			if prevLines[i] != line {
				fmt.Printf("\033[%d;1H\033[K%s", i+1, line)
			}
		}
		fmt.Printf("\033[%d;1H", len(newLines)+1)
	}

	prevLines = newLines
	time.Sleep(5 * time.Millisecond)
}

var prevLines []string

func EnterVisualMode() {
	fmt.Print(AltScreenOn)
	fmt.Print(HideCursor)
}

func ExitVisualMode() {
	fmt.Print(ShowCursor)
	fmt.Print(AltScreenOff)
}
