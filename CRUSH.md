## nona

This document provides a set of guidelines for working on the `nona` codebase.

### Commands

- **Build:** `go build -o nona .`
- **Run:** `kitty -e zsh -c 'go run .; exec zsh'`
- **Test:** `go test ./...`
- **Test Single:** `go test -run ^TestMyFunction$`
- **Lint:** `go fmt ./... && go vet ./...`

### Code Style

- **Imports:** Grouped in a single `import ()` block.
- **Formatting:** Use `go fmt` before committing.
- **Types:** Use structs for complex data models.
- **Naming:** Follow Go's standard conventions.
- **Error Handling:** Use `if err != nil` to check for errors. In `main`, use `log.Fatalf` to handle fatal errors.
- **State Machines:** Use `switch` statements on the current state to handle actions and state transitions. This pattern is used for both cell actions and model updates.
- **Testing:** Use the standard `testing` package. Use `t.Run` for table-driven tests. Each test should be isolated and self-contained, initializing its own state. Use a `bytes.Buffer` to capture and verify log output for specific actions.
