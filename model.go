package main

import (
	"fmt"
	"log"

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
	menuIndex      int
	loadedPacks    int
	totalLevels    int
	solvedLevels   int
	saveIndicators map[int]string
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
		if m.state == gameView {
			m.engine.View(*m)
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
	var s string
	if m.state != menuView {
		var title string
		if m.engine != nil {
			title = fmt.Sprintf("%s - %s by %s", m.engine.GetGameName(), m.engine.GetLevel().Name, m.engine.GetLevel().Author)
		} else {
			title = "chronical"
		}
		s = lipgloss.NewStyle().Background(lipgloss.Color("130")).Padding(0, 1).Render(title)
		s += "\n\n"
	}

	switch m.state {
	case menuView:
		s += m.viewMenuView()
	case browseView:
		s += m.viewBrowseView()
	case gameView:
		s += m.engine.View(m)
	case exportView:
		s += m.viewExportView()
	}

	return s
}
