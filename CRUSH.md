## nona

This document provides a set of guidelines for working on the `nona` codebase.

### Codebase Overview

The `nona` project is a terminal-based puzzle game framework.

- **`main.go`**: Entry point and CLI command handling using `cobra`.
- **`store.go`**: Manages the persistence of levels and saves using an SQLite database (`nona.db`).
- **`engine.go`**: Defines the core `GameEngine` interface and provides a base `Engine` struct with common game logic (grid, state, etc.).
- **Game-Specific Engines**: Each game type (e.g., `nonogram.go`) implements the `GameEngine` interface and embeds the base `Engine` to build its specific logic.

### Commands

- **Build:** `go build -o nona .`
- **Run:** `kitty -e zsh -c 'go run .; exec zsh'`
  - *Note: This command can only be run by the user.*
- **Test:** `go test ./...`
- **Test Single:** `go test -run ^TestMyFunction$`
- **Lint:** `go fmt ./... && go vet ./...`
- **Test Render:** `go run . test render --engine nonogram --initial "..." --save "..."`
  - *Note: This command is used for visual inspection of the TUI output, not for diff comparison.*
- **View Logs:** `tail -f nona.log`
  - *Note: The application writes debug output to `nona.log`.*

### Code Style

- **Imports:** Grouped in a single `import ()` block.
- **Formatting:** Use `go fmt` before committing.
- **Types:** Use structs for complex data models.
- **Naming:** Follow Go's standard conventions.
- **Error Handling:** Use `if err != nil` to check for errors. In `main`, use `log.Fatalf` to handle fatal errors.
- **State Machines:** Use `switch` statements on the current state to handle actions and state transitions. This pattern is used for both cell actions and model updates.
- **Testing:** Use the standard `testing` package. Use `t.Run` for table-driven tests. Each test should be isolated and self-contained, initializing its own state. Use a `bytes.Buffer` to capture and verify log output for specific actions.
- **Architecture**: Game-specific engines (like `NonogramEngine`) should embed the base `Engine` struct to inherit common functionality (state, grid, etc.).
- **TUI/Rendering**:
    - Use `charmbracelet/lipgloss` for all TUI rendering.
    - Style cell states (e.g., `filled`, `empty`, `cursor`) using `lipgloss.Style` with background and foreground colors rather than printing characters or drawing manual borders.
