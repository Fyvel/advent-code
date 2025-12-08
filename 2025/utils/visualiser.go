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
	BgOrange     = "\033[48;5;208m"
	BgCyan       = "\033[48;5;51m"
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
	var buf strings.Builder
	buf.WriteString(HideCursor)
	buf.WriteString(MoveCursor)

	for rowIdx, row := range grid {
		for colIdx, cell := range row {
			isActive := rowIdx == activeRow && colIdx == activeCol
			key := fmt.Sprintf("%d_%d", rowIdx, colIdx)
			isInActivePath := activePath != nil && activePath[key]

			ctx := CellRenderContext{
				Cell:           cell,
				IsActive:       isActive,
				IsInActivePath: isInActivePath,
			}

			if cellRenderer == nil {
				cellRenderer = func(ctx CellRenderContext) string {
					return ctx.Cell
				}
			}
			buf.WriteString(cellRenderer(ctx))
		}
		buf.WriteString("\n")
	}

	fmt.Print(buf.String())
	time.Sleep(5 * time.Millisecond)
}
