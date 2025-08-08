package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateGameView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		save := m.engine.GetSave()
		level := m.engine.GetLevel()
		if save.State != level.Initial {
			if err := m.store.UpsertSave(save); err != nil {
				log.Printf("event=\"save_progress_failed\" level_id=%d err=\"%v\"", level.ID, err)
			} else {
				log.Printf("event=\"save_progress_success\" level_id=%d solved=\"%v\"", level.ID, save.Solved)
			}
		}
		m.state = menuView
		m.engine = nil
		return m, nil
	case "up", "k":
		if m.engine.HasCell(m.cursorX, m.cursorY-1) {
			m.cursorY--
		}
	case "down", "j":
		if m.engine.HasCell(m.cursorX, m.cursorY+1) {
			m.cursorY++
		}
	case "left", "h":
		if m.engine.HasCell(m.cursorX-1, m.cursorY) {
			m.cursorX--
		}
	case "right", "l":
		if m.engine.HasCell(m.cursorX+1, m.cursorY) {
			m.cursorX++
		}
	case "z":
		m.engine.PrimaryAction(m.cursorX, m.cursorY)
		m.engine.Evaluate()
	case "x":
		m.engine.SecondaryAction(m.cursorX, m.cursorY)
		m.engine.Evaluate()
	case "backspace":
		m.engine.ClearCell(m.cursorX, m.cursorY)
		m.engine.Evaluate()
	}

	return m, cmd
}
