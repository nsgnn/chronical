package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("130")).Padding(0, 1)
	//define styles as variables here. No other styling should be used.
)

func (m model) View() string {
	s := titleStyle.Render("Nona")
	s += "\n\n"

	switch m.state {
	case menuView:
		s += "Welcome to Nona!\n\n"
		s += "Press 'b' to browse levels.\n"
		s += "Press 'q' to quit.\n"
	case browseView:
		if m.levels == nil {
			s += "Select a level pack:\n\n"
			for i, lp := range m.levelpacks {
				if i == m.levelPackIndex {
					s += fmt.Sprintf("> %s by %s\n", lp.Name, lp.Author)
				} else {
					s += fmt.Sprintf("  %s by %s\n", lp.Name, lp.Author)
				}
			}
		} else {
			s += fmt.Sprintf("Select a level in %s:\n\n", m.levelpacks[m.levelPackIndex].Name)
			for i, l := range m.levels {
				if i == m.levelIndex {
					s += fmt.Sprintf("> %s\n", l.Name)
				} else {
					s += fmt.Sprintf("  %s\n", l.Name)
				}
			}
		}

		s += "\nPress 'esc' to return to the menu.\n"
	case gameView:
		s += m.engine.View(m.cursorX, m.cursorY)
		s += fmt.Sprintf("\n\nCursor: (%d, %d)\n", m.cursorX, m.cursorY)
		s += "z: primary, x: secondary, backspace: clear\n"
		s += "arrow keys or hjkl to move\n"
		s += "Press 'esc' to return to the menu.\n"
	}

	return s
}
