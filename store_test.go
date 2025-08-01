package main

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestGetSaveIndicators(t *testing.T) {
	dbFile := "test_save_indicators.db"
	os.Remove(dbFile)
	store, err := NewStore(dbFile)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}
	defer store.db.Close()
	defer os.Remove(dbFile)

	// Insert a level pack for context
	levelPack := &LevelPack{Name: "Test Pack", Author: "Tester"}
	if err := store.UpsertLevelPack(levelPack); err != nil {
		t.Fatalf("failed to insert level pack: %v", err)
	}

	// Insert levels
	levels := []Level{
		{Name: "Level 1", Author: "Tester", Initial: "00", Solution: "11", Engine: "test"},
		{Name: "Level 2", Author: "Tester", Initial: "00", Solution: "11", Engine: "test"},
		{Name: "Level 3", Author: "Tester", Initial: "00", Solution: "11", Engine: "test"},
	}
	var levelIDs []int
	for i, level := range levels {
		if err := store.UpsertLevel(&levels[i], levelPack.ID); err != nil {
			t.Fatalf("failed to insert level: %v", err)
		}
		// Retrieve the full level to get its ID
		fullLevel, err := store.GetLevelByName(level.Name, levelPack.ID)
		if err != nil {
			t.Fatalf("failed to retrieve inserted level: %v", err)
		}
		levelIDs = append(levelIDs, fullLevel.ID)
	}

	// Insert saves for levels 1 and 3
	saves := []Save{
		{LevelID: levelIDs[0], State: "10", Solved: false},
		{LevelID: levelIDs[2], State: "11", Solved: true},
	}
	for _, save := range saves {
		if err := store.UpsertSave(&save); err != nil {
			t.Fatalf("failed to insert save: %v", err)
		}
	}

	// Expected indicators
	expected := map[int]string{
		levelIDs[0]: "-", // In progress
		levelIDs[1]: " ", // Not started
		levelIDs[2]: "*", // Solved
	}

	// Get indicators
	indicators, err := store.GetSaveIndicators(levelIDs)
	if err != nil {
		t.Fatalf("failed to get save indicators: %v", err)
	}

	// Compare results
	if !reflect.DeepEqual(expected, indicators) {
		t.Errorf("expected %v, got %v", expected, indicators)
	}
}

// GetLevelByName retrieves a level by its name and level pack ID.
func (s *Store) GetLevelByName(name string, levelPackID int) (*Level, error) {
	row := s.db.QueryRow(`
        SELECT id, name, author, initial_state, solution, engine
        FROM levels
        WHERE name = ? AND level_pack_id = ?;
    `, name, levelPackID)
	level := &Level{}
	err := row.Scan(&level.ID, &level.Name, &level.Author, &level.Initial, &level.Solution, &level.Engine)
	if err != nil {
		return nil, err
	}
	return level, nil
}
