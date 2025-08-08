## chronical

This document provides a set of guidelines for working on the `chronical` codebase.

### Codebase Overview

The `chronical` project is a terminal-based puzzle game framework.

- **`main.go`**: Entry point and CLI command handling using `cobra`.
- **`store.go`**: Manages the persistence of levels and saves using an SQLite database (`chronical.db`).
- **`engine.go`**: Defines the core `GameEngine` interface and provides a base `Engine` struct with common game logic (grid, state, etc.).
- **Game-Specific Engines**: Each game type (e.g., `nonogram.go`) implements the `GameEngine` interface and embeds the base `Engine` to build its specific logic.

### Commands

- **Build:** `go build -o chronical .`
- **Run:** `kitty -e zsh -c 'go run .; exec zsh'`
  - *Note: This command can only be run by the user.*
- **Test:** `go test ./...`
- **Test Single:** `go test -run ^TestMyFunction$`
- **Lint:** `go fmt ./... && go vet ./...`
- **View Logs:** `tail -f chronical.log`

### Code Style

- **Formatting:** Use `go fmt` before committing.
- **Testing:** Use the standard `testing` package with table-driven tests.
- **Architecture**: Game logic is implemented in engines that embed the base `Engine` struct.
- **TUI/Rendering**: Use `charmbracelet/lipgloss` for all TUI rendering.
