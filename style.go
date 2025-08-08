package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle      = lipgloss.NewStyle().Background(lipgloss.Color("130")).Padding(0, 1)
	subtleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	focusedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	tableTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("145"))
	blurredStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

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
