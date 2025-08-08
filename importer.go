package main

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type LevelPackYAML struct {
	Name        string  `yaml:"name"`
	Author      string  `yaml:"author"`
	Version     int     `yaml:"version"`
	Description string  `yaml:"description"`
	Levels      []Level `yaml:"levels"`
}

func (s *Store) ExportLevelPack(levelPackID int, path string) error {
	levelPack, err := s.GetLevelPack(levelPackID)
	if err != nil {
		return err
	}

	levels, err := s.GetLevelsByPack(levelPackID)
	if err != nil {
		return err
	}

	for i := range levels {
		levels[i].Initial = strings.ReplaceAll(levels[i].Initial, " ", ".")
		levels[i].Solution = strings.ReplaceAll(levels[i].Solution, " ", ".")
	}

	levelPackYAML := LevelPackYAML{
		Name:        levelPack.Name,
		Author:      levelPack.Author,
		Version:     levelPack.Version,
		Description: levelPack.Description,
		Levels:      levels,
	}

	data, err := yaml.Marshal(&levelPackYAML)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (s *Store) ImportLevelPack(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var levelPackYAML LevelPackYAML
	if err := yaml.Unmarshal(data, &levelPackYAML); err != nil {
		return err
	}

	levelPack := &LevelPack{
		Name:        levelPackYAML.Name,
		Author:      levelPackYAML.Author,
		Version:     levelPackYAML.Version,
		Description: levelPackYAML.Description,
	}

	if err := s.UpsertLevelPack(levelPack); err != nil {
		return err
	}

	for _, level := range levelPackYAML.Levels {
		level.Initial = strings.ReplaceAll(level.Initial, ".", " ")
		level.Solution = strings.ReplaceAll(level.Solution, ".", " ")
		level.SetDimensions()
		if err := s.UpsertLevel(&level, levelPack.ID); err != nil {
			return err
		}
	}

	log.Println("Successfully imported level pack:")
	log.Printf("  Name: %s\n", levelPack.Name)
	log.Printf("  Author: %s\n", levelPack.Author)
	log.Printf("  Version: %d\n", levelPack.Version)
	log.Printf("  Description: %s\n", levelPack.Description)
	log.Printf("  Levels: %d\n", len(levelPackYAML.Levels))

	return nil
}
