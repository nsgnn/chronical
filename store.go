package main

import "log"

// Store handles database access
type Store struct {
	// In a real application, this would be a database connection
	levels []Level
}

// Init initializes the store, for now with dummy data.
func (s *Store) Init() error {
	// Sample levels
	level1, err := LoadLevel(1, "Easy 4x4", "Admin", "1..4\n.2.3\n3.1.\n4..2")
	if err != nil {
		log.Printf("Error loading level 1: %v", err)
		return err
	}

	level2, err := LoadLevel(2, "Simple 3x3", "Admin", "1.3\n.2.\n3.1")
	if err != nil {
		log.Printf("Error loading level 2: %v", err)
		return err
	}

	s.levels = append(s.levels, *level1, *level2)
	return nil
}

// GetAllLevels retrieves all levels from the store.
func (s *Store) GetAllLevels() ([]Level, error) {
	return s.levels, nil
}
