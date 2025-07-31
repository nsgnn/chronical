package main

import (
	"errors"
	"log"
	"strings"
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
	l.SetDimensions()
	if l.Width <= 0 || l.Height <= 0 {
		return errors.New("level dimensions cannot be zero or negative")
	}
	return nil
}

func (l *Level) CreateSave(state string, solved bool) *Save {
	if err := l.Validate(); err != nil {
		log.Printf("event=\"invalid_level_save\" level_id=%d, err=\"%v\"", l.ID, err)
	} else {
		log.Printf("event=\"valid_level_save\" level_id=%d, state=\"%v\"", l.ID, state)
	}
	if solved {
		log.Printf("event=\"solved_level_saved\" level_id=%d", l.ID)
	}
	return &Save{
		LevelID: l.ID,
		State:   state,
		Solved:  solved,
	}
}

func (l *Level) SetDimensions() {
	if l.Initial == "" {
		l.Width = 0
		l.Height = 0
		return
	}
	lines := strings.Split(strings.TrimSpace(l.Initial), "\n")
	l.Height = len(lines)
	if l.Height > 0 {
		l.Width = len(lines[0])
	} else {
		l.Width = 0
	}
}

func (l Level) FilterValue() string { return l.Name }

func (l Level) Title() string { return l.Name }

func (l Level) ItemDescription() string { return "By " + l.Author }
