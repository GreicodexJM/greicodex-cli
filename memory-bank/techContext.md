# 📘 Tech Context — GRX CLI

## 1. Core Technology Stack
- **Language:** **Go (Golang)** was chosen for its performance, static typing, cross-compilation capabilities, and suitability for building single-binary command-line tools.
- **CLI Framework:** **`Cobra`** will be used to create a robust and POSIX-compliant CLI interface with commands, subcommands, and flags.
- **ANSI Output:** **`fatih/color`** will be used for lightweight, cross-platform colored output in the terminal to improve readability.
- **Template Management:** Go's built-in **`go:embed`** feature will be used to embed all necessary templates directly into the binary, ensuring the tool is offline-friendly.

## 2. Development and Build Setup
- **Project Structure:** The project will follow the standard Go project layout to ensure consistency and maintainability.
  ```
  grei-cli/
  ├── cmd/           # Main application entry points
  │   └── grei/
  │       └── main.go
  ├── internal/      # Private application and library code
  │   ├── core/      # Core domain logic (agnostic of frameworks)
  │   │   ├── initializer/
  │   │   └── verifier/
  │   ├── adapters/  # Implementations of ports
  │   │   ├── cli/
  │   │   ├── filesystem/
  │   │   └── reporter/
  │   └── ports/     # Interfaces for core logic
  │       ├── inbound/
  │       └── outbound/
  ├── templates/     # Embedded templates
  └── pkg/           # Public library code (if any)
  ```
- **Dependency Management:** Go Modules (`go.mod` and `go.sum`) will be used for managing project dependencies.
- **Distribution:** **`GoReleaser`** will automate the build and packaging process, creating distributables for multiple platforms, including:
  - Single binaries for Linux, macOS, and Windows.
  - `.deb` packages for Debian-based systems.
  - A Homebrew tap for macOS users.
  - `.tar.gz` archives.

## 3. Technical Constraints and Dependencies
- **Go Version:** The project will target a recent, stable version of Go (e.g., 1.21 or later) to leverage modern language features.
- **External Tools:**
  - **`gitleaks`:** The `verify` command will attempt to use `gitleaks` if it is available in the system's `PATH`. If not, it will fall back to a built-in regex-based secret scanner.
  - **Plugins:** The plugin system relies on executables named `grei-<plugin>` being present in the `PATH`.
- **Coverage Parsing:** The CLI will need to parse various test coverage report formats. It will support:
  - **JSON:** `coverage-summary.json` (from Jest).
  - **XML:** Cobertura, JaCoCo, PHPUnit, and Pytest formats.

## 4. Tool Usage Patterns
- **Template Overrides:** While templates are embedded, the CLI will allow users to override them by placing custom versions in specific directories. The lookup order will be:
  1. `.grei/templates/` (project-specific)
  2. `~/.config/grei/templates/` (user-specific)
  3. Embedded templates (default)
- **Docker Compose Handling:** The CLI will respect the convention of using `docker-compose.yml` for base service definitions (QA/PROD) and `docker-compose.override.yml` for local development overrides. This separation ensures that local development configurations do not accidentally leak into production settings.
