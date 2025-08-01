package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nona",
	Short: "A terminal-based nonogram game.",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("nona.db")
		if err != nil {
			log.Fatalf("unable to init store: %v", err)
		}

		// If the database is new, import the base levels.
		if _, err := os.Stat("nona.db"); os.IsNotExist(err) {
			if err := store.ImportLevelPack("baselevels.yaml"); err != nil {
				log.Fatalf("unable to import base levels: %v", err)
			}
		}

		m := NewModel(store)

		p := tea.NewProgram(&m)
		if _, err := p.Run(); err != nil {
			log.Fatalf("event=\"tui_failed\" err=\"%v\"", err)
		}
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a level pack to a YAML file.",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("nona.db")
		if err != nil {
			log.Fatalf("unable to init store: %v", err)
		}

		m := NewModel(store)
		m.state = exportView

		p := tea.NewProgram(&m)
		if _, err := p.Run(); err != nil {
			log.Fatalf("event=\"tui_failed\" err=\"%v\"", err)
		}
	},
}

var importCmd = &cobra.Command{
	Use:   "import [path]",
	Short: "Import a level pack from a YAML file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("nona.db")
		if err != nil {
			log.Fatalf("unable to init store: %v", err)
		}

		if err := store.ImportLevelPack(args[0]); err != nil {
			log.Fatalf("unable to import level pack: %v", err)
		}

		fmt.Printf("Level pack imported from %s\n", args[0])
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Tools for testing the game.",
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a game state for debugging.",
	Run: func(cmd *cobra.Command, args []string) {
		engineName, _ := cmd.Flags().GetString("engine")
		initialState, _ := cmd.Flags().GetString("initial")
		saveState, _ := cmd.Flags().GetString("save")
		initialState = strings.ReplaceAll(initialState, "\\n", "\n")
		saveState = strings.ReplaceAll(saveState, "\\n", "\n")
		initialState = strings.TrimSpace(initialState)
		saveState = strings.TrimSpace(saveState)

		level := Level{
			ID:       -1,
			Name:     "Test Render",
			Engine:   engineName,
			Initial:  initialState,
			Solution: saveState,
		}
		save := &Save{
			State: initialState,
		}

		var engine GameEngine
		switch engineName {
		case "nonogram":
			engine = &NonogramEngine{}
		default:
			log.Fatalf("unknown engine: %s", engineName)
		}

		game, err := engine.New(level, save)
		if err != nil {
			log.Fatalf("failed to create game: %v", err)
		}

		m := model{
			engine:  game,
			cursorX: 0,
			cursorY: 0,
		}
		fmt.Println(game.View(m))
	},
}

func init() {
	renderCmd.Flags().String("engine", "nonogram", "The engine to use for rendering.")
	renderCmd.Flags().String("initial", "", "The initial state of the grid.")
	renderCmd.Flags().String("save", "", "The save state of the grid.")

	testCmd.AddCommand(renderCmd)

	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(testCmd)
}

func main() {
	f, err := os.OpenFile("nona.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("unable to open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("event=\"root_command_failed\" err=\"%v\"", err)
	}
}
