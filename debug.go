package main

import "github.com/charmbracelet/lipgloss"

const (
	DebugPrimaryTile   = 'P'
	DebugSecondaryTile = 'S'
	DebugEmptyTile     = ' '
)

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
	return e.setCellValue(x, y, DebugPrimaryTile)
}

func (e *DebugEngine) SecondaryAction(x, y int) error {
	return e.setCellValue(x, y, DebugSecondaryTile)
}

func (e *DebugEngine) View(m model) string {
	g := e.gridView(m)
	h := e.helpView(m)
	s := lipgloss.JoinVertical(lipgloss.Left, g, h)
	return s
}

// --- Private Functions ---

func (e *DebugEngine) gridView(_ model) string {
	return "no loaded engine"
}

func (e *DebugEngine) helpView(m model) string {
	var s string
	if e.Grid[m.cursorY][m.cursorX].state != given {
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
