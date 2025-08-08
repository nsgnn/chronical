package main

// LevelPack represents a collection of levels.
type LevelPack struct {
	ID          int
	Name        string
	Author      string
	Version     int
	Description string
}

func (lp LevelPack) FilterValue() string { return lp.Name }

func (lp LevelPack) Title() string { return lp.Name }

func (lp LevelPack) ItemDescription() string { return lp.Description }
