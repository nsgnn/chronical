package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	filledTile     = '1'
	knownEmptyTile = 'X'
)

// NonogramEngine implements the GameEngine interface for nonogram puzzles.
type NonogramEngine struct {
	BaseEngine
	RowHints [][]int
	ColHints [][]int
}

func (e *NonogramEngine) New(l Level, s *Save) (GameEngine, error) {
	if _, err := e.BaseEngine.New(l, s); err != nil {
		return nil, fmt.Errorf("failed to create base engine: %w", err)
	}

	e.RowHints, e.ColHints = e.parseHints(l.Solution)

	return e, nil
}

func (e *NonogramEngine) parseHints(solution string) ([][]int, [][]int) {
	lines := strings.Split(solution, "\n")
	height := len(lines)
	width := 0
	if height > 0 {
		width = len(lines[0])
	}

	rowHints := make([][]int, height)
	for y, line := range lines {
		rowHints[y] = calculateHint(line)
	}

	colHints := make([][]int, width)
	for x := 0; x < width; x++ {
		var colBuilder strings.Builder
		for y := 0; y < height; y++ {
			colBuilder.WriteByte(lines[y][x])
		}
		colHints[x] = calculateHint(colBuilder.String())
	}

	return rowHints, colHints
}

func calculateHint(line string) []int {
	var hints []int
	count := 0
	for _, char := range line {
		if char == filledTile {
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
	return hints
}

func (e *NonogramEngine) PrimaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue(filledTile)
	e.updateSaveState()
	return nil
}

func (e *NonogramEngine) SecondaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue(knownEmptyTile)
	e.updateSaveState()
	return nil
}

func (e *NonogramEngine) EvaluateSolution() (bool, error) {
	if e.Save.Solved {
		return true, nil
	}

	sanitizedSave := strings.ReplaceAll(e.Save.State, string(knownEmptyTile), " ")
	currentRows, currentCol := e.parseHints(sanitizedSave)

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

	e.Save.Solved = true
	return true, nil
}

func (e *NonogramEngine) View(cursorX, cursorY int) string {
	var b strings.Builder

	maxRowHints := maxHintLength(e.RowHints)
	maxColHints := maxHintLength(e.ColHints)

	renderHints(&b, e.ColHints, maxRowHints, maxColHints)
	renderGrid(&b, e, cursorX, cursorY, maxRowHints)

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

func maxHintLength(hints [][]int) int {
	max := 0
	for _, hint := range hints {
		if len(hint) > max {
			max = len(hint)
		}
	}
	return max
}

func renderHints(b *strings.Builder, hints [][]int, maxRowHints, maxColHints int) {
	for i := 0; i < maxColHints; i++ {
		b.WriteString(strings.Repeat(" ", maxRowHints*3))
		for _, hint := range hints {
			if i < len(hint) {
				b.WriteString(fmt.Sprintf("  %-2d ", hint[i]))
			} else {
				b.WriteString("     ")
			}
		}
		b.WriteRune('\n')
	}
}

func renderGrid(b *strings.Builder, e *NonogramEngine, cursorX, cursorY, maxRowHints int) {
	b.WriteString(strings.Repeat("   ", maxRowHints*3))
	b.WriteString(strings.Repeat("----", e.Width()))
	b.WriteRune('\n')

	for y, row := range e.Grid {
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

		gridRowRender := lipgloss.JoinHorizontal(lipgloss.Top, rowCellRenders...)
		gridRowLines := strings.Split(gridRowRender, "\n")

		hints := e.RowHints[y]
		var hintsBuilder strings.Builder
		hintsBuilder.WriteString(strings.Repeat(" ", (maxRowHints-len(hints))*3))
		for _, hint := range hints {
			hintsBuilder.WriteString(fmt.Sprintf("%2d ", hint))
		}
		hintsString := hintsBuilder.String()
		hintsPadding := strings.Repeat(" ", maxRowHints*3)

		// A cell is 3 characters tall. The hints should be on the middle line.
		for i := range 3 {
			if i == 1 {
				b.WriteString(hintsString)
			} else {
				b.WriteString(hintsPadding)
			}
			if i < len(gridRowLines) {
				b.WriteString(gridRowLines[i])
			}
			b.WriteRune('\n')
		}
	}
}
