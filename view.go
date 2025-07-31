package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("130")).Padding(0, 1)
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

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

func (m model) View() string {
	s := titleStyle.Render("Nona Engine")
	s += "\n\n"

	switch m.state {
	case menuView:
		s += "Welcome to Nona!\n\n"
		s += subtleStyle.Render("Press 'b' to browse levels.") + "\n"
		s += subtleStyle.Render("Press 'q' to quit.") + "\n"
	case browseView:
		if m.levels == nil {
			s += "Select a level pack:\n\n"
			for i, lp := range m.levelpacks {
				if i == m.levelPackIndex {
					s += focusedStyle.Render(fmt.Sprintf("> %s by %s", lp.Name, lp.Author)) + "\n"
				} else {
					s += blurredStyle.Render(fmt.Sprintf("  %s by %s", lp.Name, lp.Author)) + "\n"
				}
			}
		} else {
			s += fmt.Sprintf("Select a level in %s:\n\n", m.levelpacks[m.levelPackIndex].Name)
			for i, l := range m.levels {
				save, err := m.store.GetSave(l.ID)
				if err != nil {
					log.Printf("event=\"no_save_file_found\" level_id=%d", l.ID)
				}

				saveStatus := "no save"
				if save != nil {
					if save.Solved {
						saveStatus = "solved"
					} else {
						saveStatus = "in progress"
					}
				}

				line := fmt.Sprintf("  %s (%s) [%s]", l.Name, l.Engine, saveStatus)
				if i == m.levelIndex {
					line = ">" + line[1:]
					s += focusedStyle.Render(line) + "\n"
				} else {
					s += blurredStyle.Render(line) + "\n"
				}
			}
		}

		s += "\n" + subtleStyle.Render("Press 'esc' to return to the menu.") + "\n"
	case gameView:
		s += m.engine.View(m.cursorX, m.cursorY)
	}

	return s
}
