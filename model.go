package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	menuView uint = iota
	browseView
	gameView
)

type model struct {
	store  *Store
	state  uint
	engine GameEngine
	//define submodels here

	levels    []Level
	listIndex int

	// game state
	cursorX int
	cursorY int
}

func NewModel(store *Store) model {
	levels, err := store.GetAllLevels()
	if err != nil {
		// In a real app, you might want to handle this more gracefully
		panic(err)
	}
	return model{
		store:  store,
		state:  menuView,
		engine: nil,
		levels: levels,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		// cmd tea.Cmd
	)
	//update submodels
	// Example:
	// m.gameboard, cmd = m.gameboard.Update(msg)
	// cmds = append(cmds, cmd)

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
				m.state = menuView
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.levels)-1 {
					m.listIndex++
				}
			case "enter":
				selectedLevel := m.levels[m.listIndex]
				var err error
				initialSave := selectedLevel.CreateSave(selectedLevel.Initial, false)
				baseEngine := &BaseEngine{}
				m.engine, err = baseEngine.New(selectedLevel, initialSave)
				if err != nil {
					panic(err)
				}
				m.state = gameView
			}
		case gameView:
			switch key {
			case "esc":
				m.state = menuView
			case "up", "k":
				m.cursorY--
			case "down", "j":
				m.cursorY++
			case "left", "h":
				m.cursorX--
			case "right", "l":
				m.cursorX++
			case "z":
				m.engine.PrimaryAction(m.cursorX, m.cursorY)
			case "x":
				m.engine.SecondaryAction(m.cursorX, m.cursorY)
			case "backspace":
				m.engine.ClearCell(m.cursorX, m.cursorY)
			}
		}
	}

	return m, tea.Batch(cmds...)
}
