package main

import (
	"errors"
	"github.com/charmbracelet/lipgloss"
	"log"
	"strings"
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
	GetSave() *Save
	GetLevel() Level
	GetGameName() string
}

func (e *BaseEngine) New(l Level, s *Save) (GameEngine, error) {
	if s == nil {
		log.Printf("event=\"EmptyLevelLoad\" level_id=%d", l.ID)
		s = l.CreateSave(l.Initial, false)
	} else {
		log.Printf("event=\"StatefulLevelLoad\" level_id=%d state=\"%v\"", l.ID, s.State)
	}

	initialLines := strings.Split(l.Initial, "\n")
	savedLines := strings.Split(s.State, "\n")

	grid := make([][]Cell, len(initialLines))
	for y, row := range initialLines {
		grid[y] = make([]Cell, len(row))
		for x, initialChar := range row {
			var cell *Cell
			savedCharRune := rune(savedLines[y][x])
			cell = NewCell(y, x, &initialChar, &savedCharRune)

			grid[y][x] = *cell
		}
	}

	e.GameName = l.Engine
	e.Level = l
	e.Save = *s
	e.Grid = grid
	return e, nil
}

func (e *BaseEngine) EvaluateSolution() (bool, error) {
	return e.Level.Solution == e.Save.State, nil
}

func (e *BaseEngine) IsValidCoordinate(x, y int) bool {
	return y >= 0 && y < len(e.Grid) && x >= 0 && x < len(e.Grid[y])
}

func (e *BaseEngine) PrimaryAction(x, y int) error {
	return errors.New("not implemented")
}

func (e *BaseEngine) SecondaryAction(x, y int) error {
	return errors.New("not implemented")
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

func (e *BaseEngine) GetSave() *Save {
	return &e.Save
}

func (e *BaseEngine) GetLevel() Level {
	return e.Level
}

func (e *BaseEngine) GetGameName() string {
	return e.GameName
}
