package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			log.Fatalf("event=\"tui_failed\" err=\"%v\"", err)
		}
	},
}

var exportCmd = &cobra.Command{
	Use:   "export [levelPackID] [path]",
	Short: "Export a level pack to a YAML file.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := NewStore("nona.db")
		if err != nil {
			log.Fatalf("unable to init store: %v", err)
		}

		levelPackID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("invalid level pack ID: %v", err)
		}

		if err := store.ExportLevelPack(levelPackID, args[1]); err != nil {
			log.Fatalf("unable to export level pack: %v", err)
		}

		fmt.Printf("Level pack %d exported to %s\n", levelPackID, args[1])
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

func init() {
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
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
