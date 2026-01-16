# AI Agent & Developer Guidelines for AWDL0 Disabler

This document serves as the primary source of truth for AI agents and developers working on the `awdl0-disabler` repository. It outlines the architectural patterns, coding standards, and operational commands required to maintain the project's integrity and style.

## 1. Environment & Build

The project is built using **Go 1.25.5** and relies on the standard Go toolchain.

### Build Commands
- **Build Binary**:
  ```bash
  go build -o build/awdl-mon cmd/awdl-mon/main.go
  ```
- **Run Directly**:
  ```bash
  go run cmd/awdl-mon/main.go
  ```
- **Download Dependencies**:
  ```bash
  go mod download
  ```

### Linting & Formatting
- **Format Code**:
  Always run `go fmt` before committing changes.
  ```bash
  go fmt ./...
  ```
- **Static Analysis**:
  Use `go vet` to catch common errors.
  ```bash
  go vet ./...
  ```
- **Modernize**:
  Check for code modernization opportunities.
  ```bash
  modernize ./...
  ```

### Testing
- **Run All Tests**:
  ```bash
  go test ./...
  ```
- **Run Tests with Verbose Output**:
  ```bash
  go test -v ./...
  ```
- **Run a Single Test**:
  To run a specific test case, use the `-run` flag with a regex matching the test name.
  *Example*: Running `TestMonitorService_Tick_DisablesWhenUp` in the services package:
  ```bash
  go test -v -run TestMonitorService_Tick_DisablesWhenUp ./internal/core/services
  ```
- **Test Coverage**:
  ```bash
  go test -cover ./...
  ```

## 2. Architecture: Hexagonal (Ports & Adapters)

This project strictly follows **Hexagonal Architecture**. Violating these boundaries breaks the core design.

### Directory Structure
- **`cmd/`**: Application entry points. Contains `main.go` which wires dependencies.
- **`internal/core/`**: The Application Core.
    - **`domain/`**: Enterprise business rules, entities, and value objects. **NO dependencies** on outer layers.
    - **`ports/`**: Interfaces that define how the core interacts with the outside world (driven ports) and how the outside world triggers the core (driving ports).
    - **`services/`**: Implementation of business logic (use cases). Depends *only* on domain and ports.
- **`internal/adapters/`**: Infrastructure and Interface implementations.
    - **`ui/`**: TUI implementation using Bubble Tea.
    - **`network/`**: Shell/Network interactions (e.g., executing `ifconfig`).
    - **`persistence/`**: Data storage implementations (memory, file, DB).

### Dependency Rule
Dependencies must point **inwards**.
- `adapters` -> depend on -> `core`
- `core` -> depends on -> `domain`
- `core` -> defines interfaces (ports)
- `adapters` -> implement interfaces

## 3. Coding Standards & Style

Adherence to these standards is mandatory.

### Control Flow: Bail First, No Else, Flat Structure
We prioritize a linear "happy path" by handling errors and edge cases immediately and returning.

**RULE: Avoid the `else` keyword.**
**RULE: Prefer "Bail First" (Early Return).**
**RULE: Avoid deep nesting. Limit if-statements to 1 level of nesting max.**

**❌ Bad (Deep Nesting / usage of Else):**
```go
func (s *Service) Process(item *Item) error {
    if item != nil {
        if item.IsValid() {
            err := s.repo.Save(item)
            if err == nil {
                return nil
            } else {
                return err
            }
        } else {
            return errors.New("invalid item")
        }
    } else {
        return errors.New("item is nil")
    }
}
```

**✅ Good (Guard Clauses / Bail First / Flat):**
```go
func (s *Service) Process(item *Item) error {
    if item == nil {
        return errors.New("item is nil")
    }

    if !item.IsValid() {
        return errors.New("invalid item")
    }

    if err := s.repo.Save(item); err != nil {
        return err
    }

    return nil
}
```

### Error Handling
- Handle errors immediately where they occur.
- Wrap errors with context if it adds value (using `fmt.Errorf("doing action: %w", err)`).
- Do not ignore errors unless explicitly safe (e.g., `_ = logger.Log(...)` in non-critical paths).

### Naming Conventions
- **Interfaces**: Name interfaces in `ports/` based on behavior (e.g., `NetworkPort`, `LoggerPort`).
- **Implementations**: Suffix with specific technology or strategy (e.g., `ShellNetworkAdapter`, `FileLogger`).
- **Variables**: Short but descriptive. `ctx` for Context, `err` for error.
- **Functions**: Verb-Noun usage (e.g., `DisableInterface`, `CheckStatus`).

### Imports
Group imports into three blocks separated by newlines:
1. Standard Library (`"fmt"`, `"os"`)
2. Project Packages (`"awdl0-disabler/internal/..."`)
3. Third-Party Libraries (`"github.com/..."`)

**Example:**
```go
import (
    "fmt"
    "time"

    "awdl0-disabler/internal/core/domain"
    "awdl0-disabler/internal/core/ports"

    tea "github.com/charmbracelet/bubbletea"
)
```

## 4. TUI Development (Bubble Tea)

The UI is built with the Charm (Bubble Tea) ecosystem.

- **Model**: Keep state minimal.
- **Update**: Handle messages (events) and return commands (`tea.Cmd`).
- **View**: Use `lipgloss` for styling. Keep view logic pure.
- **Commands**: Side effects (I/O) must be wrapped in `tea.Cmd`. Never perform I/O directly in `Update`.

## 5. Contributing Workflow

1.  **Analyze**: Understand the requirement and identify the affected hexagonal layers.
2.  **Define Port**: If a new interaction is needed, define the interface in `internal/core/ports`.
3.  **Implement Core**: Write the business logic in `internal/core/services`.
4.  **Test Core**: Write unit tests for the service using mocks.
5.  **Implement Adapter**: Create the adapter in `internal/adapters`.
6.  **Wire**: Update `cmd/awdl-mon/main.go` to inject the new adapter.
7.  **Verify**: Run `go test ./...` and `go fmt ./...`.
