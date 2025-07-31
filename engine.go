package main

import (
	"errors"
	"fmt"
	"strings"
)

// BaseEngine implements the GameEngine interface.
type BaseEngine struct {
	Name  string
	Level Level
	Save  Save
	Grid  [][]Cell
}

type GameEngine interface {
	New(l Level, s *Save) (GameEngine, error)
	EvaluateSolution() (bool, error)
	PrimaryAction(x, y int) error
	SecondaryAction(x, y int) error
	ClearCell(x, y int) error
	View() string
}

func (e *BaseEngine) New(l Level, s *Save) (GameEngine, error) {
	if s == nil {
		s = l.CreateSave(l.Initial, false)
	}

	lines := strings.Split(s.State, "\n")
	grid := make([][]Cell, len(lines))
	for y, line := range lines {
		grid[y] = make([]Cell, len(line))
		for x, r := range line {
			grid[y][x] = *NewCell(x, y, &r)
		}
	}

	return &BaseEngine{
		Name:  "DebugEngine",
		Level: l,
		Save:  *s,
		Grid:  grid,
	}, nil
}

func (e *BaseEngine) EvaluateSolution() (bool, error) {
	return e.Level.Solution == e.Save.State, nil
}

func (e *BaseEngine) isValidCoordinate(x, y int) bool {
	return y < len(e.Grid) && x < len(e.Grid[y]) && y >= 0 && x >= 0
}

func (e *BaseEngine) PrimaryAction(x, y int) error {
	if !e.isValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue('P') // Placeholder
	e.updateSaveState()
	return nil
}

func (e *BaseEngine) SecondaryAction(x, y int) error {
	if !e.isValidCoordinate(x, y) { // TODO: extract this logic to a function that accepts (x, y) and returns if it is a valid coordinate.
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue('S') // Placeholder
	e.updateSaveState()
	return nil
}

func (e *BaseEngine) ClearCell(x, y int) error {
	if !e.isValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].Clear()
	e.updateSaveState()
	return nil
}

func (e *BaseEngine) View() string {
	return fmt.Sprintf("%s\nsolved: %t", e.Save.State, e.Save.Solved)
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
