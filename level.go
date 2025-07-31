package main

import (
	"errors"
	"log"
)

type Level struct {
	ID       int    `yaml:"id" json:"id"`
	Name     string `yaml:"name" json:"name"`
	Author   string `yaml:"author" json:"author"`
	Initial  string `yaml:"initial" json:"initial"`
	Solution string `yaml:"solution" json:"solution"`
	Engine   string `yaml:"engine" json:"engine"`
	Width    int    `yaml:"width" json:"width"`
	Height   int    `yaml:"height" json:"height"`
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
