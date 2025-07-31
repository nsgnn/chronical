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
			BorderForeground(lipgloss.Color("242")) // Gray
	filledStyle  = cellStyle
	invalidStyle = cellStyle.
			BorderForeground(lipgloss.Color("196")) // Red
)

func (m model) View() string {
	var title string
	if m.engine != nil {
		title = fmt.Sprintf("%s - %s", m.engine.GetGameName(), m.engine.GetLevel().Name)
	} else {
		title = "Nona Engine"
	}
	s := titleStyle.Render(title)
	s += "\n\n"

	switch m.state {
	case menuView:
		loadedPacks := len(m.levelpacks)
		totalLevels, err := m.store.CountLevels()
		if err != nil {
			log.Printf("could not count levels %v", err)
		}
		solvedLevels, err := m.store.CountSolvedLevels()
		if err != nil {
			log.Printf("could not count solved levels %v", err)
		}

		var solvedPercentage float64
		if totalLevels > 0 {
			solvedPercentage = (float64(solvedLevels) / float64(totalLevels)) * 100
		}

		stats := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render(fmt.Sprintf(
				"Loaded Packs: %d\nTotal Levels: %d\nSolved: %d/%d (%.2f%%)",
				loadedPacks,
				totalLevels,
				solvedLevels,
				totalLevels,
				solvedPercentage,
			))

		s += stats + "\n\n"
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

				saveIndicator := " "
				if save != nil {
					if save.Solved {
						saveIndicator = "*"
					} else {
						saveIndicator = "-"
					}
				}

				line := fmt.Sprintf("  %s\t(%s)\t%s", l.Name, l.Engine, saveIndicator)
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
