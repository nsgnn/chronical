package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chronical",
	Short: "A cli-based puzzle engine supporting a variety of modes and plain text levelpacks.",
	Long: `Chronical is a terminal-based puzzle game engine that supports a variety of puzzle types.
It uses a plain text, YAML-based format for creating and sharing level packs, making it easy for anyone to create their own puzzles.`,
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("chronical.db")
		if err != nil {
			log.Fatalf("unable to init store: %v", err)
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
	Short: "Export a level pack to a YAML file. This is useful for sharing level packs with others.",
	Long: `Export a level pack to a YAML file. This is useful for sharing level packs with others.
The exported file can be imported by other users using the import command.`,
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("chronical.db")
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
	Short: "Import a level pack from a YAML file. This is useful for playing level packs created by others.",
	Long: `Import a level pack from a YAML file. This is useful for playing level packs created by others.
The imported file will be added to your library of level packs.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("chronical.db")
		if err != nil {
			log.Fatalf("unable to init store: %v", err)
		}

		if err := store.ImportLevelPack(args[0]); err != nil {
			log.Fatalf("unable to import level pack: %v", err)
		}

		fmt.Printf("Level pack imported from %s\n", args[0])
	},
}

var testImportCmd = &cobra.Command{
	Use:   "import [path]",
	Short: "Test importing and exporting a level pack to ensure that the process is working correctly.",
	Long: `Test importing and exporting a level pack to ensure that the process is working correctly.
This command is useful for developers who want to test the import/export functionality.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("chronical.db")
		if err != nil {
			log.Fatalf("unable to init store: %v", err)
		}

		// Import the level pack
		if err := store.ImportLevelPack(args[0]); err != nil {
			log.Fatalf("unable to import level pack: %v", err)
		}

		// Get the last inserted level pack which is what we just inserted
		rows, err := store.db.Query(`
			SELECT id, name
			FROM level_packs
			ORDER BY id DESC
			LIMIT 1;
		`)
		if err != nil {
			log.Fatalf("unable to get last level pack: %v", err)
		}
		defer rows.Close()

		var id int
		var name string
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatalf("unable to scan last level pack: %v", err)
			}
		}

		// Export the level pack to a temporary file
		tmpfile, err := os.CreateTemp("", "exported-level-pack-*.yaml")
		if err != nil {
			log.Fatalf("unable to create temporary file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		if err := store.ExportLevelPack(id, tmpfile.Name()); err != nil {
			log.Fatalf("unable to export level pack: %v", err)
		}

		// Compare the two files
		original, err := os.ReadFile(args[0])
		if err != nil {
			log.Fatalf("unable to read original file: %v", err)
		}
		exported, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			log.Fatalf("unable to read exported file: %v", err)
		}

		var originalYAML, exportedYAML LevelPackYAML
		if err := yaml.Unmarshal(original, &originalYAML); err != nil {
			log.Fatalf("unable to unmarshal original file: %v", err)
		}
		if err := yaml.Unmarshal(exported, &exportedYAML); err != nil {
			log.Fatalf("unable to unmarshal exported file: %v", err)
		}

		originalNormalized, err := yaml.Marshal(&originalYAML)
		if err != nil {
			log.Fatalf("unable to marshal original file: %v", err)
		}
		exportedNormalized, err := yaml.Marshal(&exportedYAML)
		if err != nil {
			log.Fatalf("unable to marshal exported file: %v", err)
		}

		if string(originalNormalized) != string(exportedNormalized) {
			log.Fatalf("exported file does not match original file")
		}

		fmt.Println("import/export test passed")
	},
}

var testCmd = &cobra.Command{
	Use:    "test",
	Short:  "Tools for testing the game.",
	Hidden: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logToStdout, _ := cmd.Flags().GetBool("log-stdout")
		if logToStdout {
			log.SetOutput(os.Stdout)
		}
	},
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a game state for debugging. This is useful for testing new game engines.",
	Long: `Render a game state for debugging. This is useful for testing new game engines.
You can use this command to see how a level will look without having to play it.`,
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
	testCmd.AddCommand(testImportCmd)
	testImportCmd.Flags().BoolP("help", "h", false, "Help message for the test import command")
	testCmd.PersistentFlags().Bool("log-stdout", false, "Write logs to stdout instead of a file.")

	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(testCmd)
}

func main() {
	f, err := os.OpenFile("chronical.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("unable to open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("event=\"root_command_failed\" err=\"%v\"", err)
	}
}
