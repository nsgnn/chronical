package main

import (
	"errors"
	"log"
)

type Level struct {
	ID       int
	Name     string
	Author   string
	Initial  string
	Solution string
	Engine   string
}

func (l *Level) Validate() error {
	if l.ID < 0 {
		return errors.New("level id cannot be negative")
	}

	if l.Solution == "" {
		return errors.New("solution cannot be empty")
	}

	return nil
}

func NewLevel(id int, name string, author string, solution string) (*Level, error) {
	l := &Level{
		ID:       id,
		Name:     name,
		Author:   author,
		Solution: solution,
		Engine:   "debug", // for now, we will always want the debug engine.
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
