package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateMenuView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < 2 {
			m.menuIndex++
		}
	case "enter":
		switch m.menuIndex {
		case 0:
			m.state = browseView
		case 1:
			m.state = exportView
		case 2:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) viewMenuView() string {
	title := `
  ▌       ▘    ▜ 
▛▘▛▌▛▘▛▌▛▌▌▛▘▀▌▐ 
▙▖▌▌▌ ▙▌▌▌▌▙▖█▌▐▖
`
	var s string
	s += lipgloss.NewStyle().Margin(1).Foreground(lipgloss.Color("130")).Render(title)

	buttons := []string{"Browse", "Export", "Quit"}
	for i, button := range buttons {
		style := lipgloss.NewStyle().Padding(1, 2)
		if i == m.menuIndex {
			style = style.Foreground(lipgloss.Color("205")).Bold(true)
			s += style.Render("> " + button)
		} else {
			s += style.Render("  " + button)
		}
		s += "\n"
	}

	var solvedPercentage float64
	if m.totalLevels > 0 {
		solvedPercentage = (float64(m.solvedLevels) / float64(m.totalLevels)) * 100
	}

	stats := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Padding(1, 2).
		Render(fmt.Sprintf(
			"Loaded: %d Pack(s) %d Levels\nSolved: %d/%d (%.2f%%)",
			m.loadedPacks,
			m.totalLevels,
			m.solvedLevels,
			m.totalLevels,
			solvedPercentage,
		))

	return lipgloss.JoinHorizontal(lipgloss.Bottom, s, stats)
}
