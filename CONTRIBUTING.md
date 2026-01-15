# Contributing to AWDL0 Disabler

Thank you for your interest in contributing! This project uses **Hexagonal Architecture** (Ports & Adapters) to ensure modularity and testability.

## 🛠 Development Setup

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/awdl0-disabler.git
    cd awdl0-disabler
    ```

2.  **Install dependencies**:
    ```bash
    go mod download
    ```

3.  **Run Tests**:
    ```bash
    go test ./...
    ```

## 📐 Architecture Guidelines

When adding features, please adhere to the Hexagonal Architecture pattern:

1.  **Core Domain (`internal/core`)**:
    *   This layer must **not** depend on any external libraries or adapters.
    *   Define logic in `services/`.
    *   Define interfaces in `ports/`.

2.  **Adapters (`internal/adapters`)**:
    *   Implement the interfaces defined in `ports/`.
    *   Keep implementation details (like CLI libraries or File I/O) contained here.

3.  **Dependency Injection**:
    *   Wire everything together in `cmd/awdl-mon/main.go`.

## 🧪 Testing

*   **Unit Tests**: Required for any new logic in `internal/core/services`.
*   **Mocks**: Use the mock implementations in `internal/core/services/mocks_test.go` to test services without side effects.

## 📝 Style Guide

*   Use `go fmt` before committing.
*   Follow standard Go idioms (Effective Go).
*   Keep commit messages clear and concise.
