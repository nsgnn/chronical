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
	BaseEngine
	RowHints [][]int
	ColHints [][]int
}

func (e *NonogramEngine) New(l Level, s *Save) (GameEngine, error) {
	_, err := e.BaseEngine.New(l, s)
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

// View returns a string representation of the game state, including hints.
func (e *NonogramEngine) View(cursorX, cursorY int) string {
	var b strings.Builder

	// Get max hint lengths for formatting
	maxRowHints := 0
	for _, hints := range e.RowHints {
		if len(hints) > maxRowHints {
			maxRowHints = len(hints)
		}
	}
	maxColHints := 0
	for _, hints := range e.ColHints {
		if len(hints) > maxColHints {
			maxColHints = len(hints)
		}
	}

	// Draw column hints
	for i := 0; i < maxColHints; i++ {
		b.WriteString(strings.Repeat(" ", maxRowHints*3)) // Spacer
		for _, hints := range e.ColHints {
			if i < len(hints) {
				b.WriteString(fmt.Sprintf("  %-2d ", hints[i]))
			} else {
				b.WriteString("     ")
			}
		}
		b.WriteRune('\n')
	}

	// Draw separator
	b.WriteString(strings.Repeat("   ", maxRowHints*3))
	b.WriteString(strings.Repeat("----", e.Width()))
	b.WriteRune('\n')

	// Draw grid and row hints
	for y, row := range e.Grid {
		// Render all cells for the current row
		var rowCellRenders []string
		for x, cell := range row {
			var style lipgloss.Style
			switch cell.state {
			case given:
				style = givenStyle
			case filled:
				style = filledStyle
			case invalid:
				style = invalidStyle
			default:
				style = cellStyle
			}

			if x == cursorX && y == cursorY {
				style = cursorStyle
			}
			rowCellRenders = append(rowCellRenders, style.Render(string(cell.value)))
		}

		// Join them horizontally to get a multi-line string for the grid row
		gridRowRender := lipgloss.JoinHorizontal(lipgloss.Top, rowCellRenders...)
		gridRowLines := strings.Split(gridRowRender, "\n")

		// Prepare hints string for the middle line
		hints := e.RowHints[y]
		var hintsBuilder strings.Builder
		hintsBuilder.WriteString(strings.Repeat(" ", (maxRowHints-len(hints))*3))
		for _, hint := range hints {
			hintsBuilder.WriteString(fmt.Sprintf("%2d ", hint))
		}
		hintsString := hintsBuilder.String()

		// Prepare padding for hints for top and bottom lines
		hintsPadding := strings.Repeat(" ", maxRowHints*3)

		// A cell is 3 characters tall. The hints should be on the middle line.
		for i := range 3 {
			if i == 1 { // Middle line for hints
				b.WriteString(hintsString)
			} else {
				b.WriteString(hintsPadding)
			}
			// b.WriteString("|")
			if i < len(gridRowLines) {
				b.WriteString(gridRowLines[i])
			}
			b.WriteRune('\n')
		}
	}

	s := b.String()
	if e.Grid[cursorY][cursorX].state != given {
		s += "\nz: Toggle, x: Flag, backspace: clear\n"
	} else {
		s += "\n\n"
	}
	s += "arrow keys or hjkl to move\n"
	s += "Press 'esc' to return to the menu.\n"
	if e.Save.Solved {
		s += "Congrats!\n"
	}
	return s
}
