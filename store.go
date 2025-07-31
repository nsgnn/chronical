package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Store handles all database operations.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store and initializes the database connection.
func NewStore(dataSourceName string) (*Store, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	store := &Store{db: db}
	if err := store.Migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

// Migrate creates the necessary database tables if they don't already exist.
func (s *Store) Migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS level_packs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			author TEXT,
			version INTEGER NOT NULL DEFAULT 1,
			description TEXT
		);
	`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS levels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level_pack_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			author TEXT,
			initial_state TEXT NOT NULL,
			solution TEXT NOT NULL,
			engine TEXT NOT NULL,
			FOREIGN KEY (level_pack_id) REFERENCES level_packs(id),
			UNIQUE(level_pack_id, name)
		);
	`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS saves (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level_id INTEGER NOT NULL UNIQUE,
			state TEXT NOT NULL,
			solved BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (level_id) REFERENCES levels(id)
		);
	`)
	return err
}

// UpsertLevelPack inserts or updates a level pack.
func (s *Store) UpsertLevelPack(pack *LevelPack) error {
	row := s.db.QueryRow(`
		INSERT INTO level_packs (name, author, version, description)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			author = excluded.author,
			version = excluded.version,
			description = excluded.description
		RETURNING id;
	`, pack.Name, pack.Author, pack.Version, pack.Description)
	return row.Scan(&pack.ID)
}

// GetLevelPack retrieves a level pack by its ID.
func (s *Store) GetLevelPack(id int) (*LevelPack, error) {
	log.Printf("event=\"get_level_pack\" id=%d", id)
	row := s.db.QueryRow(`
		SELECT id, name, author, version, description
		FROM level_packs
		WHERE id = ?;
	`, id)
	pack := &LevelPack{}
	err := row.Scan(&pack.ID, &pack.Name, &pack.Author, &pack.Version, &pack.Description)
	if err != nil {
		return nil, err
	}
	log.Printf("event=\"found_level_pack\" name=\"%s\"", pack.Name)
	return pack, nil
}

// GetLevelPacks retrieves all level packs.
func (s *Store) GetAllLevelPacks() ([]LevelPack, error) {
	log.Println("event=\"get_all_level_packs\"")
	rows, err := s.db.Query(`
		SELECT id, name, author, version, description
		FROM level_packs;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packs []LevelPack
	for rows.Next() {
		pack := LevelPack{}
		err := rows.Scan(&pack.ID, &pack.Name, &pack.Author, &pack.Version, &pack.Description)
		if err != nil {
			return nil, err
		}
		packs = append(packs, pack)
	}
	log.Printf("event=\"found_level_packs\" count=%d", len(packs))
	return packs, nil
}

// UpsertLevel inserts or updates a level.
func (s *Store) UpsertLevel(level *Level, levelPackID int) error {
	_, err := s.db.Exec(`
		INSERT INTO levels (level_pack_id, name, author, initial_state, solution, engine)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(level_pack_id, name) DO UPDATE SET
			author = excluded.author,
			initial_state = excluded.initial_state,
			solution = excluded.solution,
			engine = excluded.engine;
	`, levelPackID, level.Name, level.Author, level.Initial, level.Solution, level.Engine)
	return err
}

// GetLevel retrieves a level by its ID.
func (s *Store) GetLevel(id int) (*Level, error) {
	log.Printf("event=\"get_level\" id=%d", id)
	row := s.db.QueryRow(`
		SELECT id, name, author, initial_state, solution, engine
		FROM levels
		WHERE id = ?;
	`, id)
	level := &Level{}
	err := row.Scan(&level.ID, &level.Name, &level.Author, &level.Initial, &level.Solution, &level.Engine)
	if err != nil {
		return nil, err
	}
	log.Printf("event=\"found_level\" name=\"%s\"", level.Name)
	return level, nil
}

// GetLevelsByPack retrieves all levels for a given level pack.
func (s *Store) GetLevelsByPack(levelPackID int) ([]Level, error) {
	log.Printf("event=\"get_levels_by_pack\" level_pack_id=%d", levelPackID)
	rows, err := s.db.Query(`
		SELECT id, name, author, initial_state, solution, engine
		FROM levels
		WHERE level_pack_id = ?;
	`, levelPackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []Level
	for rows.Next() {
		level := Level{}
		err := rows.Scan(&level.ID, &level.Name, &level.Author, &level.Initial, &level.Solution, &level.Engine)
		if err != nil {
			return nil, err
		}
		levels = append(levels, level)
	}
	log.Printf("event=\"found_levels_for_pack\" count=%d level_pack_id=%d", len(levels), levelPackID)
	return levels, nil
}

// UpsertSave inserts or updates a save.
func (s *Store) UpsertSave(save *Save) error {
	level, err := s.GetLevel(save.LevelID)
	if err != nil {
		return err
	}
	if level.Initial == save.State {
		log.Printf("event=\"delete_save_on_upsert\" level_id=%d", save.LevelID)
		return s.DeleteSave(save.LevelID)
	}
	_, err = s.db.Exec(`
		INSERT INTO saves (level_id, state, solved, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(level_id) DO UPDATE SET
			state = excluded.state,
			solved = excluded.solved,
			updated_at = CURRENT_TIMESTAMP;
	`, save.LevelID, save.State, save.Solved)
	return err
}

// GetSave retrieves a save by its level ID.
func (s *Store) GetSave(levelID int) (*Save, error) {
	log.Printf("event=\"get_save\" level_id=%d", levelID)
	row := s.db.QueryRow(`
		SELECT level_id, state, solved, created_at, updated_at
		FROM saves
		WHERE level_id = ?;
	`, levelID)
	save := &Save{}
	err := row.Scan(&save.LevelID, &save.State, &save.Solved, &save.CreatedAt, &save.UpdatedAt)
	if err != nil {
		return nil, err
	}
	log.Printf("event=\"found_save\" level_id=%d", levelID)
	return save, nil
}

// DeleteSave deletes a save by its level ID.
func (s *Store) DeleteSave(levelID int) error {
	_, err := s.db.Exec(`
		DELETE FROM saves
		WHERE level_id = ?;
	`, levelID)
	return err
}

// GetAllLevels is added to satisfy the model.go dependency.
func (s *Store) GetAllLevels() ([]Level, error) {
	log.Println("event=\"get_all_levels\"")
	rows, err := s.db.Query(`
		SELECT id, name, author, initial_state, solution, engine
		FROM levels;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []Level
	for rows.Next() {
		level := Level{}
		err := rows.Scan(&level.ID, &level.Name, &level.Author, &level.Initial, &level.Solution, &level.Engine)
		if err != nil {
			return nil, err
		}
		levels = append(levels, level)
	}
	log.Printf("event=\"found_levels\" count=%d", len(levels))
	return levels, nil
}
