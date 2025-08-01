package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	menuView uint = iota
	browseView
	gameView
	exportView
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type model struct {
	store          *Store
	state          uint
	engine         GameEngine
	cursorX        int
	cursorY        int
	levelpacks     []LevelPack
	levels         []Level
	levelPackIndex int
	levelIndex     int
	loadedPacks    int
	totalLevels    int
	solvedLevels   int
	saveIndicators map[int]string
	viewport       viewport.Model
}

func NewModel(store *Store) model {
	levelpacks, err := store.GetAllLevelPacks()
	if err != nil {
		log.Fatalf("could not get all level packs: %v", err)
	}

	loadedPacks := len(levelpacks)
	totalLevels, err := store.CountLevels()
	if err != nil {
		log.Printf("could not count levels %v", err)
	}
	solvedLevels, err := store.CountSolvedLevels()
	if err != nil {
		log.Printf("could not count solved levels %v", err)
	}

	return model{
		store:          store,
		state:          menuView,
		engine:         nil,
		levelpacks:     levelpacks,
		loadedPacks:    loadedPacks,
		totalLevels:    totalLevels,
		solvedLevels:   solvedLevels,
		saveIndicators: make(map[int]string),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 3
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - headerHeight
		if m.state == gameView {
			m.viewport.SetContent(m.engine.View(*m))
		}
		return m, nil
	case errMsg:
		log.Printf("error: %v", msg)
		return m, tea.Quit
	case tea.KeyMsg:
		switch m.state {
		case menuView:
			return m.updateMenuView(msg)
		case browseView:
			return m.updateBrowseView(msg)
		case gameView:
			return m.updateGameView(msg)
		case exportView:
			return m.updateExportView(msg)
		}
	}
	return m, nil
}

func (m model) View() string {
	var title string
	if m.engine != nil {
		title = fmt.Sprintf("%s - %s by %s", m.engine.GetGameName(), m.engine.GetLevel().Name, m.engine.GetLevel().Author)
	} else {
		title = "Nona Engine"
	}
	s := lipgloss.NewStyle().Background(lipgloss.Color("130")).Padding(0, 1).Render(title)
	s += "\n\n"

	switch m.state {
	case menuView:
		s += m.viewMenuView()
	case browseView:
		s += m.viewBrowseView()
	case gameView:
		s += m.viewport.View()
	case exportView:
		s += m.viewExportView()
	}

	return s
}

func (m *model) updateMenuView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		return m, tea.Quit
	case "b":
		m.state = browseView
	case "e":
		m.state = exportView
	}
	return m, nil
}

func (m model) viewMenuView() string {
	var solvedPercentage float64
	if m.totalLevels > 0 {
		solvedPercentage = (float64(m.solvedLevels) / float64(m.totalLevels)) * 100
	}

	stats := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(fmt.Sprintf(
			"Loaded Packs: %d\nTotal Levels: %d\nSolved: %d/%d (%.2f%%)",
			m.loadedPacks,
			m.totalLevels,
			m.solvedLevels,
			m.totalLevels,
			solvedPercentage,
		))

	s := stats + "\n\n"
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press 'b' to browse levels.") + "\n"
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press 'q' to quit.") + "\n"
	return s
}

func (m *model) updateBrowseView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
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
			m.viewport.SetContent(m.engine.View(*m))
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
		return m, nil
	case "up", "k":
		if m.engine.IsValidCoordinate(m.cursorX, m.cursorY-1) {
			m.cursorY--
		}
	case "down", "j":
		if m.engine.IsValidCoordinate(m.cursorX, m.cursorY+1) {
			m.cursorY++
		}
	case "left", "h":
		if m.engine.IsValidCoordinate(m.cursorX-1, m.cursorY) {
			m.cursorX--
		}
	case "right", "l":
		if m.engine.IsValidCoordinate(m.cursorX+1, m.cursorY) {
			m.cursorX++
		}
	case "z":
		m.engine.PrimaryAction(m.cursorX, m.cursorY)
		m.engine.EvaluateSolution()
	case "x":
		m.engine.SecondaryAction(m.cursorX, m.cursorY)
		m.engine.EvaluateSolution()
	case "backspace":
		m.engine.ClearCell(m.cursorX, m.cursorY)
		m.engine.EvaluateSolution()
	}

	m.viewport.SetContent(m.engine.View(*m))
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

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
