// This file implements the Nonogram game logic.
//
// Nonogram is a picture logic puzzle in which cells in a grid must be
// colored or left blank according to numbers at the side of the grid to
// reveal a hidden picture. In this puzzle type, the numbers are a form of
// discrete tomography that measures how many unbroken lines of filled-in
// squares there are in any given row or column.
//
// For example, a clue of "4 8 3" would mean there are sets of four, eight,
// and three filled squares, in that order, with at least one blank square
// between successive sets.
//
// The primary action will color a cell with the FilledTile rune.
// The secondary action will mark a cell as a known empty tile with the KnownEmptyTile rune.
// The puzzle is evaluated as a solve if the save and solution tomography matches for both rows and columns.

package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	FilledTile     = rune('1')
	KnownEmptyTile = rune('X')
	EmptyTile      = rune(' ')
)

var (
	renderStyles = map[rune]lipgloss.Style{
		FilledTile:     lipgloss.NewStyle().Background(lipgloss.Color("255")), // White background
		KnownEmptyTile: lipgloss.NewStyle().Foreground(lipgloss.Color("245")), // Light grey foreground
		EmptyTile:      lipgloss.NewStyle().Foreground(lipgloss.Color("250")), // Dim grey foreground
	}
	highlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("205")) // Magenta background for the cursor
	hintStyle      = lipgloss.NewStyle()
	renderRunes    = map[rune]string{
		FilledTile:     " ",
		KnownEmptyTile: " X",
		EmptyTile:      " ·",
	}
	cellWidth = 2
)

type NonogramEngine struct {
	Engine
	rowHints      [][]int
	colHints      [][]int
	hintRowWidth  int
	hintColHeight int
}

func (e *NonogramEngine) New(l Level, s *Save) (GameEngine, error) {
	_, err := e.Engine.New(l, s)
	if err != nil {
		return nil, err
	}

	e.rowHints, e.colHints = generateTomography(l.Solution)

	for _, hints := range e.colHints {
		if len(hints) > e.hintColHeight {
			e.hintColHeight = len(hints)
		}
	}

	for _, hints := range e.rowHints {
		if len(hints) > e.hintRowWidth {
			e.hintRowWidth = len(hints)
		}
	}
	e.hintRowWidth *= 3 // Each hint is typically 3 characters " 4 " or "10 "

	return e, nil
}

func (e *NonogramEngine) PrimaryAction(x, y int) error {
	return e.setCellValue(x, y, FilledTile)
}

func (e *NonogramEngine) SecondaryAction(x, y int) error {
	return e.setCellValue(x, y, KnownEmptyTile)
}

func (e *NonogramEngine) Evaluate() (bool, error) {
	r, c := generateTomography(e.Save.State)

	result := cmpTomo(r, e.rowHints) && cmpTomo(c, e.colHints)
	return result, nil
}

func (e *NonogramEngine) View(m model) string {
	g := e.gridView(m)
	r := e.rowHintView()
	c := e.colHintView()
	h := e.helpView(m)

	spacer := hintStyle.Width(e.hintRowWidth).Height(e.hintColHeight).Render("")

	s1 := lipgloss.JoinHorizontal(lipgloss.Bottom, spacer, c)
	s2 := lipgloss.JoinHorizontal(lipgloss.Top, r, g)
	s := lipgloss.JoinVertical(lipgloss.Left, s1, s2, h)

	return s
}

// --- Private Functions ---

func tileView(c Cell, h bool) string {
	s, ok := renderStyles[c.value]
	if !ok {
		//TODO: structured log out the unknown tile.
		s = renderStyles[EmptyTile]
	}
	if h {
		s = highlightStyle
	}
	r, ok := renderRunes[c.value]
	if !ok {
		//TODO structured log out the unknown tile.
		r = renderRunes[EmptyTile]
	}
	return s.Width(cellWidth).Render(r)
}

func (e *NonogramEngine) gridView(m model) string {
	var rows []string
	for y, row := range e.Grid {
		var rowBuider []string
		for x, cell := range row {
			highlighted := x == m.cursorX && y == m.cursorY
			cell := tileView(cell, highlighted)
			rowBuider = append(rowBuider, cell)

		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, rowBuider...)
		rows = append(rows, row)
	}
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return grid
}

func (e *NonogramEngine) colHintView() string {

	var columns []string
	for i := 0; i < e.hintColHeight; i++ {
		var row []string
		for _, hints := range e.colHints {
			i := len(hints) - e.hintColHeight + i
			var b string
			if i >= 0 {
				b = fmt.Sprintf("%d", hints[i])
			} else {
				b = " "
			}
			row = append(row, hintStyle.Width(cellWidth).Align(lipgloss.Center).Render(b))
		}
		columns = append(columns, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}
	s := lipgloss.JoinVertical(lipgloss.Left, columns...)
	return s
}

func (e *NonogramEngine) rowHintView() string {
	var rows []string
	for _, hints := range e.rowHints {
		var b []string
		for _, h := range hints {
			b = append(b, fmt.Sprintf("%2d", h))
		}
		s := strings.Join(b, " ")
		rows = append(rows, hintStyle.Width(e.hintRowWidth).Align(lipgloss.Right).Render(s))
	}
	s := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return s
}

func (e *NonogramEngine) helpView(m model) string {
	help := "\n"
	if e.Grid[m.cursorY][m.cursorX].state != given {
		help += "z: Toggle\tx: Mark Empty\tbackspace: clear\n"
	} else {
		help += "\n"
	}
	help += "arrow keys or hjkl to move\n"
	help += "Press 'esc' to return to the menu.\n"
	if e.Save.Solved {
		help += "Congrats!\n"
	}
	return help
}

func generateTomography(state string) ([][]int, [][]int) {
	s := strings.Split(state, "\n")
	h := len(s)
	w := 0
	if h > 0 {
		for _, row := range s {
			if len(row) > w {
				w = len(row)
			}
		}
	}

	rowHints := make([][]int, h)
	for y, row := range s {
		c := 0
		var rowHint []int
		for _, r := range row {
			if r == FilledTile {
				c++
			} else {
				if c > 0 {
					rowHint = append(rowHint, c)
				}
				c = 0
			}
		}
		if c > 0 {
			rowHint = append(rowHint, c)
		}
		if len(rowHint) == 0 {
			rowHint = append(rowHint, 0)
		}
		rowHints[y] = rowHint
	}

	colHints := make([][]int, w)
	for x := 0; x < w; x++ {
		c := 0
		var colHint []int
		for y := 0; y < h; y++ {
			r := EmptyTile
			if x < len(s[y]) {
				r = rune(s[y][x])
			}

			if r == FilledTile {
				c++
			} else {
				if c > 0 {
					colHint = append(colHint, c)
				}
				c = 0
			}
		}
		if c > 0 {
			colHint = append(colHint, c)
		}
		if len(colHint) == 0 {
			colHint = append(colHint, 0)
		}
		colHints[x] = colHint
	}

	return rowHints, colHints
}

// cmpTomo compares two sets of tomography hints for equality.
func cmpTomo(a [][]int, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, r := range a {
		if len(r) != len(b[i]) {
			return false
		}
		for j, v := range r {
			if v != b[i][j] {
				return false
			}
		}
	}
	return true
}
