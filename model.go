package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	menuView uint = iota
	browseView
	gameView
)

type model struct {
	store *Store
	state uint

	// game state
	engine  GameEngine
	cursorX int
	cursorY int

	// menu state
	levelpacks     []LevelPack
	levels         []Level
	levelPackIndex int
	levelIndex     int
}

func NewModel(store *Store) model {
	levelpacks, err := store.GetAllLevelPacks()
	if err != nil {
		// In a real app, you might want to handle this more gracefully
		panic(err)
	}

	return model{
		store:      store,
		state:      menuView,
		engine:     nil,
		levelpacks: levelpacks,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case menuView:
			switch key {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			case "b":
				m.state = browseView
			}
		case browseView:
			switch key {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "esc":
				// If we're looking at levels, go back to level packs.
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
					// We're looking at level packs, so switch to levels.
					selectedPack := m.levelpacks[m.levelPackIndex]
					levels, err := m.store.GetLevelsByPack(selectedPack.ID)
					if err != nil {
						panic(err)
					}
					m.levels = levels
					m.levelIndex = 0
				} else {
					// We're looking at levels, so switch to the game.
					selectedLevel := m.levels[m.levelIndex]

					// Try to get an existing save.
					save, err := m.store.GetSave(selectedLevel.ID)
					if err != nil {
						log.Println("No save file found, creating a new one.")
					}

					baseEngine := &BaseEngine{}
					m.engine, err = baseEngine.New(selectedLevel, save)
					if err != nil {
						panic(err)
					}

					m.state = gameView
					m.cursorX = 0
					m.cursorY = 0
				}
			}
		case gameView:
			switch key {
			case "esc":
				baseEngine := m.engine.(*BaseEngine)
				if baseEngine.Save.State == baseEngine.Level.Initial {
					if err := m.store.DeleteSave(baseEngine.Level.ID); err != nil {
						log.Printf("failed to delete save: %v", err)
					}
				} else {
					if err := m.store.UpsertSave(&baseEngine.Save); err != nil {
						log.Printf("failed to save progress: %v", err)
					}
				}
				m.state = menuView
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
			case "x":
				m.engine.SecondaryAction(m.cursorX, m.cursorY)
			case "backspace":
				m.engine.ClearCell(m.cursorX, m.cursorY)
			}
		}
	}

	return m, nil
}
