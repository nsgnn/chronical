package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateExportView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.state = menuView
	case "up", "k":
		if m.levelPackIndex > 0 {
			m.levelPackIndex--
		}
	case "down", "j":
		if m.levelPackIndex < len(m.levelpacks)-1 {
			m.levelPackIndex++
		}
	case "enter":
		selectedPack := m.levelpacks[m.levelPackIndex]
		err := m.store.ExportLevelPack(selectedPack.ID, selectedPack.Name+".yaml")
		if err != nil {
			return m, func() tea.Msg { return errMsg{err} }
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m model) viewExportView() string {
	s := "Select a level pack to export:\n\n"
	for i, lp := range m.levelpacks {
		if i == m.levelPackIndex {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(fmt.Sprintf("> %s by %s", lp.Name, lp.Author)) + "\n"
		} else {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Render(fmt.Sprintf("  %s by %s", lp.Name, lp.Author)) + "\n"
		}
	}
	s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press 'enter' to export the selected level pack.") + "\n"
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press 'esc' to return to the menu.") + "\n"
	return s
}
