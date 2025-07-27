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
		s += "Select a level:\n\n"
		for i, level := range m.levels {
			cursor := " "
			if m.listIndex == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s by %s\n", cursor, level.name, level.author)
		}
		s += "\nPress 'enter' to play a level.\n"
		s += "Press 'esc' to return to the menu.\n"
	case gameView:
		s += "Game view is not yet implemented.\n\n"
		s += "Press 'esc' to return to the menu.\n"
	}

	return s
}
