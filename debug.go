package main

import "errors"

// DebugEngine implements the GameEngine interface for debugging.
type DebugEngine struct {
	Engine
}

func (e *DebugEngine) New(l Level, s *Save) (GameEngine, error) {
	_, err := e.Engine.New(l, s)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *DebugEngine) PrimaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue('P') // Placeholder
	e.updateSaveState()
	return nil
}

func (e *DebugEngine) SecondaryAction(x, y int) error {
	if !e.IsValidCoordinate(x, y) {
		return errors.New("coordinates out of bounds")
	}
	e.Grid[y][x].EnterValue('S') // Placeholder
	e.updateSaveState()
	return nil
}
