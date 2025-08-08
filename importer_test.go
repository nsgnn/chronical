package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestImporter(t *testing.T) {
	staticYAML := `
name: Test Pack
author: Crush
version: 1
description: A test level pack
levels:
  - id: 1
    name: Test Level 1
    author: Crush
    initial: |
      .1.
      1.1
      .1.
    solution: |
      111
      111
      111
    engine: nonogram
    width: 3
    height: 3
`
	db, err := NewStore("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to create in-memory store: %v", err)
	}

	tmpfile, err := os.CreateTemp("", "test.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(staticYAML)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	t.Run("ImportAndExport", func(t *testing.T) {
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		log.SetOutput(w)

		if err := db.ImportLevelPack(tmpfile.Name()); err != nil {
			t.Fatalf("failed to import level pack: %v", err)
		}

		w.Close()
		os.Stderr = oldStderr
		log.SetOutput(os.Stderr)

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r)
		if err != nil {
			t.Fatalf("failed to read from pipe: %v", err)
		}
		output := strings.TrimSpace(buf.String())

		if !strings.Contains(output, "Successfully imported level pack") {
			t.Errorf("expected log to contain success message, got %q", output)
		}

		levelPacks, err := db.GetAllLevelPacks()
		if err != nil {
			t.Fatalf("failed to get level packs: %v", err)
		}
		if len(levelPacks) != 1 {
			t.Fatalf("expected 1 level pack, got %d", len(levelPacks))
		}
		levelPack := levelPacks[0]
		if levelPack.Name != "Test Pack" {
			t.Errorf("expected level pack name to be 'Test Pack', got %q", levelPack.Name)
		}

		levels, err := db.GetLevelsByPack(levelPack.ID)
		if err != nil {
			t.Fatalf("failed to get levels by pack: %v", err)
		}
		if len(levels) != 1 {
			t.Errorf("expected 1 level, got %d", len(levels))
		}

		exportPath := tmpfile.Name() + ".exported"
		defer os.Remove(exportPath)

		if err := db.ExportLevelPack(levelPack.ID, exportPath); err != nil {
			t.Fatalf("failed to export level pack: %v", err)
		}

		data, err := os.ReadFile(exportPath)
		if err != nil {
			t.Fatalf("failed to read exported file: %v", err)
		}

		var original, exported LevelPackYAML
		if err := yaml.Unmarshal([]byte(staticYAML), &original); err != nil {
			t.Fatalf("failed to unmarshal original YAML: %v", err)
		}
		if err := yaml.Unmarshal(data, &exported); err != nil {
			t.Fatalf("failed to unmarshal exported YAML: %v", err)
		}

		if original.Name != exported.Name {
			t.Errorf("expected name %q, got %q", original.Name, exported.Name)
		}
		if original.Author != exported.Author {
			t.Errorf("expected author %q, got %q", original.Author, exported.Author)
		}
		if len(original.Levels) != len(exported.Levels) {
			t.Fatalf("expected %d levels, got %d", len(original.Levels), len(exported.Levels))
		}
		if original.Levels[0].Name != exported.Levels[0].Name {
			t.Errorf("expected level name %q, got %q", original.Levels[0].Name, exported.Levels[0].Name)
		}
	})

	t.Run("UpdatePack", func(t *testing.T) {
		// First, import the original pack.
		if err := db.ImportLevelPack(tmpfile.Name()); err != nil {
			t.Fatalf("failed to import level pack: %v", err)
		}

		// Now, create an updated YAML for the same pack.
		updatedYAML := `
name: Test Pack
author: Crush
version: 2
description: An updated test level pack
levels:
  - id: 1
    name: Test Level 1
    author: Crush
    initial: |
      .1.
      1.1
      .1.
    solution: |
      111
      111
      111
    engine: nonogram
    width: 3
    height: 3
  - id: 2
    name: Test Level 2
    author: Crush
    initial: |
      1.
      .1
    solution: |
      1.
      .1
    engine: nonogram
    width: 2
    height: 2
`
		if err := os.WriteFile(tmpfile.Name(), []byte(updatedYAML), 0644); err != nil {
			t.Fatalf("failed to write updated YAML: %v", err)
		}

		// Import the updated pack.
		if err := db.ImportLevelPack(tmpfile.Name()); err != nil {
			t.Fatalf("failed to import updated level pack: %v", err)
		}

		// Verify the pack has been updated.
		levelPacks, err := db.GetAllLevelPacks()
		if err != nil {
			t.Fatalf("failed to get level packs: %v", err)
		}
		if len(levelPacks) != 1 {
			t.Fatalf("expected 1 level pack, got %d", len(levelPacks))
		}
		levelPack := levelPacks[0]
		if levelPack.Version != 2 {
			t.Errorf("expected version 2, got %d", levelPack.Version)
		}
		if levelPack.Description != "An updated test level pack" {
			t.Errorf("expected description to be updated, got %q", levelPack.Description)
		}

		// Verify the levels have been updated.
		levels, err := db.GetLevelsByPack(levelPack.ID)
		if err != nil {
			t.Fatalf("failed to get levels by pack: %v", err)
		}
		if len(levels) != 2 {
			t.Fatalf("expected 2 levels, got %d", len(levels))
		}
	})
}
