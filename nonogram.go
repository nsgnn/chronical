package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	FilledTile     = rune('1')
	KnownEmptyTile = rune('X')
)

// NonogramEngine implements the GameEngine interface for nonogram puzzles.
type NonogramEngine struct {
	Engine
	RowHints [][]int
	ColHints [][]int
}

func (e *NonogramEngine) New(l Level, s *Save) (GameEngine, error) {
	_, err := e.Engine.New(l, s)
	if err != nil {
		return nil, err
	}

	e.RowHints, e.ColHints = e.generateHints(l.Solution)

	return e, nil
}

// generateHints parses the solution string and returns the row and column hints.
func (e *NonogramEngine) generateHints(solution string) ([][]int, [][]int) {
	lines := strings.Split(solution, "\n")
	height := len(lines)
	width := 0
	if height > 0 {
		width = len(lines[0])
	}

	rowHints := make([][]int, height)
	for y := range height {
		count := 0
		var hints []int
		for x := 0; x < width; x++ {
			if rune(lines[y][x]) == FilledTile {
				count++
			} else {
				if count > 0 {
					hints = append(hints, count)
				}
				count = 0
			}
		}
		if count > 0 {
			hints = append(hints, count)
		}
		if len(hints) == 0 {
			hints = append(hints, 0)
		}
		rowHints[y] = hints
	}

	colHints := make([][]int, width)
	for x := 0; x < width; x++ {
		count := 0
		var hints []int
		for y := range height {
			if rune(lines[y][x]) == FilledTile {
				count++
			} else {
				if count > 0 {
					hints = append(hints, count)
				}
				count = 0
			}
		}
		if count > 0 {
			hints = append(hints, count)
		}
		if len(hints) == 0 {
			hints = append(hints, 0)
		}
		colHints[x] = hints
	}

	return rowHints, colHints
}

func (e *NonogramEngine) PrimaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue(FilledTile)
	e.updateSaveState()
	return nil
}

func (e *NonogramEngine) SecondaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue(KnownEmptyTile)
	e.updateSaveState()
	return nil
}

func (e *NonogramEngine) EvaluateSolution() (bool, error) {
	sanitized_known_empty_save := strings.ReplaceAll(e.Save.State, string(KnownEmptyTile), " ")
	currentRows, currentCol := e.generateHints(sanitized_known_empty_save)

	// Compare row hints
	if len(currentRows) != len(e.RowHints) {
		return false, nil
	}
	for i, row := range currentRows {
		if len(row) != len(e.RowHints[i]) {
			return false, nil
		}
		for j, val := range row {
			if val != e.RowHints[i][j] {
				return false, nil
			}
		}
	}

	// Compare column hints
	if len(currentCol) != len(e.ColHints) {
		return false, nil
	}
	for i, col := range currentCol {
		if len(col) != len(e.ColHints[i]) {
			return false, nil
		}
		for j, val := range col {
			if val != e.ColHints[i][j] {
				return false, nil
			}
		}
	}

	return true, nil
}

func (e *NonogramEngine) View(cursorX, cursorY int) string {
	// --- Styles ---
	styleFilled := lipgloss.NewStyle().Background(lipgloss.Color("255"))     // White background
	styleKnownEmpty := lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // Light grey foreground
	styleEmpty := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))      // Dim grey foreground
	styleCursor := lipgloss.NewStyle().Background(lipgloss.Color("205"))     // Magenta background for the cursor

	// --- Characters ---
	charFilled := " "
	charKnownEmpty := "X"
	charEmpty := "·"

	// Each grid cell is two characters wide for better proportions.
	cellWidth := 2

	// --- Grid Rendering ---
	var gridRows []string
	for y, row := range e.Grid {
		var rowCells []string
		for x, cell := range row {
			var style lipgloss.Style
			var char string

			switch cell.value {
			case FilledTile:
				style = styleFilled
				char = charFilled
			case KnownEmptyTile:
				style = styleKnownEmpty
				char = charKnownEmpty
			default:
				style = styleEmpty
				char = charEmpty
			}

			// Pad character to cellWidth
			paddedChar := lipgloss.NewStyle().Width(cellWidth).Render(char)

			if x == cursorX && y == cursorY {
				rowCells = append(rowCells, styleCursor.Render(paddedChar))
			} else {
				rowCells = append(rowCells, style.Render(paddedChar))
			}
		}
		gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, rowCells...))
	}
	gridView := lipgloss.JoinVertical(lipgloss.Left, gridRows...)

	// --- Row Hints Rendering ---
	maxRowHintsLen := 0
	for _, hints := range e.RowHints {
		if len(hints) > maxRowHintsLen {
			maxRowHintsLen = len(hints)
		}
	}
	// Each hint number takes about 3 spaces " 1 "
	rowHintsWidth := maxRowHintsLen * 3

	var rowHintViews []string
	for _, hints := range e.RowHints {
		var hintStrings []string
		for _, h := range hints {
			hintStrings = append(hintStrings, fmt.Sprintf("%2d", h))
		}
		s := strings.Join(hintStrings, " ")
		rowHintViews = append(rowHintViews, lipgloss.NewStyle().Width(rowHintsWidth).Align(lipgloss.Right).Render(s))
	}
	rowHintsView := lipgloss.JoinVertical(lipgloss.Left, rowHintViews...)

	// --- Column Hints Rendering ---
	maxColHintsLen := 0
	for _, hints := range e.ColHints {
		if len(hints) > maxColHintsLen {
			maxColHintsLen = len(hints)
		}
	}

	var colHintStrs []string
	for i := 0; i < maxColHintsLen; i++ {
		var row []string
		for _, hints := range e.ColHints {
			hintIdx := len(hints) - maxColHintsLen + i
			var hintStr string
			if hintIdx >= 0 {
				hintStr = fmt.Sprintf("%d", hints[hintIdx])
			} else {
				hintStr = " "
			}
			row = append(row, lipgloss.NewStyle().Width(cellWidth).Align(lipgloss.Center).Render(hintStr))
		}
		colHintStrs = append(colHintStrs, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}
	colHintsView := lipgloss.JoinVertical(lipgloss.Left, colHintStrs...)

	// --- Assembly ---
	spacer := lipgloss.NewStyle().Width(rowHintsWidth).Height(maxColHintsLen).Render("")
	topView := lipgloss.JoinHorizontal(lipgloss.Bottom, spacer, colHintsView)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, rowHintsView, gridView)
	finalView := lipgloss.JoinVertical(lipgloss.Left, topView, mainView)

	// --- Help Text ---
	help := "\n"
	if e.Grid[cursorY][cursorX].state != given {
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
