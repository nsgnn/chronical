package main

import (
	"errors"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	cellStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			Padding(0, 1)
	cursorStyle = cellStyle.
			Border(lipgloss.ThickBorder(), true)
	givenStyle = cellStyle.
			BorderForeground(lipgloss.Color("242")) // Gray for given cells
	filledStyle  = cellStyle
	invalidStyle = cellStyle.
			BorderForeground(lipgloss.Color("196")) // Red border for invalid
)

// BaseEngine implements the GameEngine interface.
type BaseEngine struct {
	GameName string
	Level    Level
	Save     Save
	Grid     [][]Cell
}

type GameEngine interface {
	New(l Level, s *Save) (GameEngine, error)
	EvaluateSolution() (bool, error)
	PrimaryAction(x, y int) error
	SecondaryAction(x, y int) error
	ClearCell(x, y int) error
	View(cursorX, cursorY int) string
	Width() int
	Height() int
	IsValidCoordinate(x, y int) bool
}

func (e *BaseEngine) New(l Level, s *Save) (GameEngine, error) {
	if s == nil {
		s = l.CreateSave(l.Initial, false)
	}

	initialLines := strings.Split(l.Initial, "\n")
	savedLines := strings.Split(s.State, "\n")

	grid := make([][]Cell, len(initialLines))
	for y, row := range initialLines {
		grid[y] = make([]Cell, len(row))
		for x, initialChar := range row {
			if initialChar != '.' {
				grid[y][x] = *NewCell(x, y, &initialChar)
				continue
			}

			savedChar := rune(savedLines[y][x])
			cell := NewCell(x, y, nil)
			if savedChar != ' ' && savedChar != '.' {
				cell.EnterValue(savedChar)
			}
			grid[y][x] = *cell
		}
	}

	return &BaseEngine{
		GameName: l.Engine,
		Level:    l,
		Save:     *s,
		Grid:     grid,
	}, nil
}

func (e *BaseEngine) EvaluateSolution() (bool, error) {
	return e.Level.Solution == e.Save.State, nil
}

func (e *BaseEngine) IsValidCoordinate(x, y int) bool {
	return y >= 0 && y < len(e.Grid) && x >= 0 && x < len(e.Grid[y])
}

func (e *BaseEngine) PrimaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue('P') // Placeholder
	e.updateSaveState()
	return nil
}

func (e *BaseEngine) SecondaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue('S') // Placeholder
	e.updateSaveState()
	return nil
}

func (e *BaseEngine) ClearCell(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].Clear()
	e.updateSaveState()
	return nil
}

func (e *BaseEngine) View(cursorX, cursorY int) string {
	var rows []string
	for y, row := range e.Grid {
		var rowStrings []string
		for x, cell := range row {
			var style lipgloss.Style
			switch cell.state {
			case given:
				style = givenStyle
			case filled:
				style = filledStyle
			case invalid:
				style = invalidStyle
			default: // empty
				style = cellStyle
			}

			if x == cursorX && y == cursorY {
				style = cursorStyle
			}

			rowStrings = append(rowStrings, style.Render(cell.View()))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, rowStrings...))
	}
	s := lipgloss.JoinVertical(lipgloss.Left, rows...)
	if e.Grid[cursorY][cursorX].state != given {
		s += "\nz: primary, x: secondary, backspace: clear\n"
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

func (e *BaseEngine) Width() int {
	return e.Level.Width
}

func (e *BaseEngine) Height() int {
	return e.Level.Height
}

func (e *BaseEngine) updateSaveState() {
	var builder strings.Builder
	for y, row := range e.Grid {
		for _, cell := range row {
			builder.WriteRune(cell.value)
		}
		if y < len(e.Grid)-1 {
			builder.WriteRune('\n')
		}
	}
	e.Save.State = builder.String()
	solved, err := e.EvaluateSolution()
	if err != nil {
		e.Save.Solved = false
	} else {
		e.Save.Solved = solved
	}
}
