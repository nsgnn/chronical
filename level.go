package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Level struct {
	id       int
	name     string
	author   string
	solution string
}

type Save struct {
	level  int
	state  string
	solved bool
}

func (l *Level) validate() error {
	if l.id < 0 {
		return errors.New("level id cannot be negative")
	}

	if l.solution == "" {
		return errors.New("solution cannot be empty")
	}

	// Normalize newlines and remove a single optional trailing newline
	// to handle text files that end with a newline character.
	normalizedSolution := strings.ReplaceAll(l.solution, "\r\n", "\n")
	processedSolution := strings.TrimSuffix(normalizedSolution, "\n")

	rows := strings.Split(processedSolution, "\n")

	if len(rows) <= 1 {
		return errors.New("solution y-dimension is too small")
	}

	rowLength := len(rows[0])
	if rowLength <= 1 {
		return errors.New("solution x-dimension is too small")
	}

	for i, row := range rows {
		if len(row) != rowLength {
			return fmt.Errorf("solution row %d has length %d, expected %d", i+1, len(row), rowLength)
		}
	}

	if processedSolution != l.solution {
		return errors.New("solution format errors")
	}

	return nil
}

func LoadLevel(id int, name string, author string, solution string) (*Level, error) {
	l := &Level{
		id:       id,
		name:     name,
		author:   author,
		solution: solution,
	}

	if err := l.validate(); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *Level) CreateSave(state string) *Save {
	if err := l.validate(); err != nil {
		log.Println("creating save for invalid level. red flag")
	}

	return &Save{
		level:  l.id,
		state:  state,
		solved: state == l.solution,
	}
}
