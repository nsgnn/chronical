package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateBrowseView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.levels != nil {
			m.levels = nil
			m.levelIndex = 0
			return m, nil
		}
		m.state = menuView
	case "up", "k":
		if m.levels == nil {
			if m.levelPackIndex > 0 {
				m.levelPackIndex--
			}
		} else {
			if m.levelIndex > 0 {
				m.levelIndex--
			}
		}
	case "down", "j":
		if m.levels == nil {
			if m.levelPackIndex < len(m.levelpacks)-1 {
				m.levelPackIndex++
			}
		} else {
			if m.levelIndex < len(m.levels)-1 {
				m.levelIndex++
			}
		}
	case "enter":
		if m.levels == nil {
			selectedPack := m.levelpacks[m.levelPackIndex]
			levels, err := m.store.GetLevelsByPack(selectedPack.ID)
			if err != nil {
				return m, func() tea.Msg { return errMsg{err} }
			}
			m.levels = levels
			m.levelIndex = 0

			var levelIDs []int
			for _, level := range levels {
				levelIDs = append(levelIDs, level.ID)
			}
			indicators, err := m.store.GetSaveIndicators(levelIDs)
			if err != nil {
				return m, func() tea.Msg { return errMsg{err} }
			}
			m.saveIndicators = indicators
		} else {
			selectedLevel := m.levels[m.levelIndex]

			save, err := m.store.GetSave(selectedLevel.ID)
			if err != nil {
				log.Printf("event=\"no_save_file_found\" level_id=%d", selectedLevel.ID)
			}

			var engine GameEngine
			switch selectedLevel.Engine {
			case "nonogram":
				engine, err = new(NonogramEngine).New(selectedLevel, save)
			default:
				log.Printf("event=\"engine_not_found\" level_engine=\"%v\"", selectedLevel.Engine)
				selectedLevel.Engine = "fallback"
				engine, err = new(DebugEngine).New(selectedLevel, save)
			}
			if err != nil {
				return m, func() tea.Msg { return errMsg{err} }
			}
			log.Printf("event=\"created_engine\" engine_type=\"%v\"", engine.GetGameName())
			m.engine = engine

			m.state = gameView
			m.cursorX = 0
			m.cursorY = 0
			m.levels = nil
		}
	}
	return m, nil
}

func (m model) viewBrowseView() string {
	var s string
	if m.levels == nil {
		s += "Select a level pack:\n\n"
		for i, lp := range m.levelpacks {
			if i == m.levelPackIndex {
				s += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("> %s by %s", lp.Name, lp.Author)) + "\n"
			} else {
				s += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render(fmt.Sprintf("  %s by %s", lp.Name, lp.Author)) + "\n"
			}
		}
	} else {
		s += fmt.Sprintf("Select a level in %s:\n\n", m.levelpacks[m.levelPackIndex].Name)
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("145")).Render(fmt.Sprintf("  %-24s\t(%s)\t%s", "Level Name", "Game Mode", "Save")) + "\n"
		for i, l := range m.levels {
			saveIndicator := m.saveIndicators[l.ID]

			line := fmt.Sprintf("  %-24s\t(%s)\t%s", l.Name, l.Engine, saveIndicator)
			if i == m.levelIndex {
				line = ">" + line[1:]
				s += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(line) + "\n"
			} else {
				s += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render(line) + "\n"
			}
		}
	}

	s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press 'esc' to return to the menu.") + "\n"
	return s
}
