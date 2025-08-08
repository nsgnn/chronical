package main

import (
	"errors"
	"log"
	"strings"
)

// GameEngine defines the interface for a game engine.
type GameEngine interface {
	New(l Level, s *Save) (GameEngine, error)
	Evaluate() (bool, error)
	PrimaryAction(x, y int) error
	SecondaryAction(x, y int) error
	setCellValue(x, y int, value rune) error
	ClearCell(x, y int) error
	View(m model) string
	GetWidth() int
	GetHeight() int
	HasCell(x, y int) bool
	GetSave() *Save
	GetLevel() Level
	GetGameName() string
}

// Engine implements the GameEngine interface.
type Engine struct {
	GameName string
	Level    Level
	Save     Save
	Grid     [][]Cell
}

func (e *Engine) New(l Level, s *Save) (GameEngine, error) {
	if s == nil {
		log.Printf("event=\"EmptyLevelLoad\" level_id=%d", l.ID)
		s = l.CreateSave(l.Initial, false)
	} else {
		log.Printf("event=\"StatefulLevelLoad\" level_id=%d state=\"%v\"", l.ID, s.State)
	}

	iRows := strings.Split(l.Initial, "\n")
	sRow := strings.Split(s.State, "\n")

	grid := make([][]Cell, len(iRows))
	for y, row := range iRows {
		grid[y] = make([]Cell, len(row))
		for x, initialChar := range row {
			var cell *Cell
			savedCharRune := rune(sRow[y][x])
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

func (e *Engine) Evaluate() (bool, error) {
	return e.Level.Solution == e.Save.State, nil
}

func (e *Engine) PrimaryAction(x, y int) error {
	return errors.New("not implemented")
}

func (e *Engine) SecondaryAction(x, y int) error {
	return errors.New("not implemented")
}

func (e *Engine) ClearCell(x, y int) error {
	if !e.HasCell(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].Clear()
	e.updateSaveState()
	return nil
}

func (e *Engine) View(m model) string {
	return ""
}

func (e *Engine) GetWidth() int {
	return len(e.Grid[0])
}

func (e *Engine) GetHeight() int {
	return len(e.Grid)
}

func (e *Engine) HasCell(x, y int) bool {
	return y >= 0 && y < len(e.Grid) && x >= 0 && x < len(e.Grid[y])
}

func (e *Engine) GetSave() *Save {
	return &e.Save
}

func (e *Engine) GetLevel() Level {
	return e.Level
}

func (e *Engine) GetGameName() string {
	return e.GameName
}

func (e *Engine) setCellValue(x, y int, value rune) error {
	if !e.HasCell(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue(value)
	e.updateSaveState()
	return nil
}

func (e *Engine) updateSaveState() {
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
	solved, err := e.Evaluate()
	if err != nil {
		e.Save.Solved = false
	} else {
		e.Save.Solved = solved
	}
}
