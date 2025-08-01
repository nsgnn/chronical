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

type NonogramEngine struct {
	Engine
	solutionRowTomography    [][]int
	solutionColumnTomography [][]int
}

func (e *NonogramEngine) New(l Level, s *Save) (GameEngine, error) {
	_, err := e.Engine.New(l, s)
	if err != nil {
		return nil, err
	}

	e.solutionRowTomography, e.solutionColumnTomography = e.generateTomography(l.Solution)

	return e, nil
}

func (e *NonogramEngine) generateTomography(state string) ([][]int, [][]int) {
	s := strings.Split(state, "\n")
	h := len(s)
	w := 0
	if h > 0 {
		w = len(s[0])
	}

	rowHints := make([][]int, h)
	for y := range h {
		c := 0
		var rowHint []int
		for x := 0; x < w; x++ {
			r := rune(s[y][x])
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
		for y := range h {
			r := rune(s[y][x])
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

func (e *NonogramEngine) PrimaryAction(x, y int) error {
	return e.SetCellValue(x, y, FilledTile)
}

func (e *NonogramEngine) SecondaryAction(x, y int) error {
	return e.SetCellValue(x, y, KnownEmptyTile)
}

func (e *NonogramEngine) EvaluateSolution() (bool, error) {
	saveRowTomo, saveColTomo := e.generateTomography(e.Save.State)

	result := cmpTomo(saveRowTomo, e.solutionRowTomography) && cmpTomo(saveColTomo, e.solutionColumnTomography)
	return result, nil
}

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

func (e *NonogramEngine) View(m model) string {

	styleMap := map[rune]lipgloss.Style{
		FilledTile:     lipgloss.NewStyle().Background(lipgloss.Color("255")), // White background
		KnownEmptyTile: lipgloss.NewStyle().Foreground(lipgloss.Color("245")), // Light grey foreground
		EmptyTile:      lipgloss.NewStyle().Foreground(lipgloss.Color("250")), // Dim grey foreground
	}
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color("205")) // Magenta background for the cursor
	hintStyle := lipgloss.NewStyle()
	charMap := map[rune]string{
		FilledTile:     " ",
		KnownEmptyTile: "X",
		EmptyTile:      "·",
	}
	// Each grid cell is two characters wide for better proportions.
	cellWidth := 2

	// --- Grid Rendering ---
	var rows []string
	for y, row := range e.Grid {
		var rowBuider []string
		for x, cell := range row {
			renderStyle, ok := styleMap[cell.value]
			if !ok {
				renderStyle = styleMap[EmptyTile]
			}
			renderRune, ok := charMap[cell.value]
			if !ok {
				renderRune = charMap[EmptyTile]
			}

			if x == m.cursorX && y == m.cursorY {
				renderStyle = cursorStyle
			}
			rowBuider = append(rowBuider, renderStyle.Width(cellWidth).Render(renderRune))

		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, rowBuider...)
		rows = append(rows, row)
	}
	grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// --- Row Hints Rendering ---
	maxRowHintsLen := 0
	for _, hints := range e.solutionRowTomography {
		if len(hints) > maxRowHintsLen {
			maxRowHintsLen = len(hints)
		}
	}
	// Each hint number takes about 3 spaces " 1 "
	rowHintsWidth := maxRowHintsLen * 3

	var rowHintViews []string
	for _, hints := range e.solutionRowTomography {
		var hintStrings []string
		for _, h := range hints {
			hintStrings = append(hintStrings, fmt.Sprintf("%2d", h))
		}
		s := strings.Join(hintStrings, " ")
		rowHintViews = append(rowHintViews, hintStyle.Width(rowHintsWidth).Align(lipgloss.Right).Render(s))
	}
	rowHintsView := lipgloss.JoinVertical(lipgloss.Left, rowHintViews...)

	// --- Column Hints Rendering ---
	maxColHintsLen := 0
	for _, hints := range e.solutionColumnTomography {
		if len(hints) > maxColHintsLen {
			maxColHintsLen = len(hints)
		}
	}

	var colHintStrs []string
	for i := 0; i < maxColHintsLen; i++ {
		var row []string
		for _, hints := range e.solutionColumnTomography {
			hintIdx := len(hints) - maxColHintsLen + i
			var hintStr string
			if hintIdx >= 0 {
				hintStr = fmt.Sprintf("%d", hints[hintIdx])
			} else {
				hintStr = " "
			}
			row = append(row, hintStyle.Width(cellWidth).Align(lipgloss.Center).Render(hintStr))
		}
		colHintStrs = append(colHintStrs, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}
	colHintsView := lipgloss.JoinVertical(lipgloss.Left, colHintStrs...)

	// --- Assembly ---
	spacer := hintStyle.Width(rowHintsWidth).Height(maxColHintsLen).Render("")
	topView := lipgloss.JoinHorizontal(lipgloss.Bottom, spacer, colHintsView)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, rowHintsView, grid)
	finalView := lipgloss.JoinVertical(lipgloss.Left, topView, mainView)

	// --- Help Text ---
	help := "\n"
	if e.Grid[m.cursorY][m.cursorX].state != given {
		help += "z: Toggle, x: Flag, backspace: clear\n"
	} else {
		help += "\n"
	}
	help += "arrow keys or hjkl to move\n"
	help += "Press 'esc' to return to the menu.\n"
	if e.Save.Solved {
		help += "Congrats!\n"
	}

	return finalView + help
}
