# Contributing to AWDL0 Disabler

Thank you for your interest in contributing! This project uses **Hexagonal Architecture** (Ports & Adapters) to ensure modularity and testability.

## 🚀 Getting Started

To contribute to this project, please follow these steps:

1.  **Fork the repository**: Click the "Fork" button at the top right of the GitHub page.
2.  **Clone your fork**:
    ```bash
    git clone https://github.com/your-username/awdl0-disabler.git
    cd awdl0-disabler
    ```
3.  **Add the upstream remote**:
    ```bash
    git remote add upstream https://github.com/anderson-oki/awdl0-disabler.git
    ```
4.  **Create a feature branch**:
    ```bash
    git checkout -b feature/your-feature-name
    ```

## 🛠 Development Setup

1.  **Install dependencies**:
    ```bash
    go mod download
    ```

2.  **Run Tests**:
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

## 📝 Pull Request Process

1.  **Update documentation**: If you've changed any functionality, update the `README.md` or other relevant documentation.
2.  **Run tests**: Ensure all tests pass before submitting.
3.  **Format code**: Use `go fmt ./...` to ensure consistent styling.
4.  **Submit Pull Request**: Open a Pull Request against the `master` branch of the upstream repository. Provide a clear description of your changes.

## 📝 Style Guide

*   Follow standard Go idioms (Effective Go).
*   Keep commit messages clear and concise.
