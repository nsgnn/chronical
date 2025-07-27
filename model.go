package main

import tea "github.com/charmbracelet/bubbletea"

const (
	menuView uint = iota
	browseView
	gameView
)

type model struct {
	store *Store
	state uint
	//define submodels here

	levels    []Level
	listIndex int
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
				// TODO: Start the game with the selected level
				m.state = gameView
			}
		case gameView:
			switch key {
			case "esc":
				m.state = menuView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
