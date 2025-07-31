package main

import (
	"errors"
	"log"
	"strings"
)

type Level struct {
	ID       int
	Name     string
	Author   string
	Initial  string
	Solution string
	Engine   string
	Width    int
	Height   int
}

func (l *Level) Validate() error {
	if l.ID < 0 {
		return errors.New("level id cannot be negative")
	}

	if l.Solution == "" {
		return errors.New("solution cannot be empty")
	}

	if l.Width <= 0 || l.Height <= 0 {
		return errors.New("level dimensions cannot be zero or negative")
	}

	return nil
}

func NewLevel(id int, name string, author string, initial string, solution string) (*Level, error) {
	lines := strings.Split(solution, "\n")
	height := len(lines)
	width := 0
	if height > 0 {
		width = len(lines[0])
	}

	l := &Level{
		ID:       id,
		Name:     name,
		Author:   author,
		Initial:  initial,
		Solution: solution,
		Engine:   "debug", // for now, we will always want the debug engine.
		Width:    width,
		Height:   height,
	}

	if err := l.Validate(); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *Level) CreateSave(state string, solved bool) *Save {
	if err := l.Validate(); err != nil {
		log.Println("creating save for invalid level. red flag")
	}

	return &Save{
		LevelID: l.ID,
		State:   state,
		Solved:  solved,
	}
}

func (l Level) FilterValue() string { return l.Name }

func (l Level) Title() string { return l.Name }

func (l Level) ItemDescription() string { return "By " + l.Author }
