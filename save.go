package main

import "time"

type Save struct {
	LevelID   int
	State     string
	Solved    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
