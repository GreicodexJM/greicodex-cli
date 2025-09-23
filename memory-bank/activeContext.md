# 📘 Active Context — GRX CLI

## 1. Current Work Focus
The primary focus is on **Phase 1: Core CLI Implementation**. This involves building the foundational commands and features of the GRX CLI.

**Immediate Goals:**
1.  **Set up the Go project structure:** Create the directory layout as defined in `techContext.md`.
2.  **Implement the `grei init` command:** This is the first command to be built. It will initialize a new project with the standard Greicodex structure and templates.
3.  **Embed initial templates:** Create a basic set of templates (e.g., `README.md`, `.gitignore`) and embed them into the binary using `go:embed`.

## 2. Recent Changes & Decisions
- **Decision:** The `memory-bank` has been initialized with the core documentation (`projectbrief.md`, `productContext.md`, `systemPatterns.md`, `techContext.md`). This provides a solid foundation for the project's goals, architecture, and technical stack.
- **Decision:** The project will be built using Go and the Cobra framework, following a Hexagonal Architecture. This decision was made to ensure the CLI is performant, maintainable, and easy to test.

## 3. Next Steps
- **Initialize Go module:** Run `go mod init` to create the `go.mod` file.
- **Create the main application entry point:** Set up `cmd/grei/main.go`.
- **Implement the root command:** Use Cobra to create the main `grei` command.
- **Add the `init` subcommand:** Implement the `grei init` command, including the logic to copy embedded templates to the target directory.

## 4. Important Patterns & Preferences
- **Spanish-first:** All user-facing output (command descriptions, help text, messages) should be in Spanish.
- **Offline-first:** The CLI must be fully functional without an internet connection. All essential resources, like templates, must be embedded.
- **Clean Code and SOLID:** Adhere to the principles outlined in the `.clinerules` to ensure the codebase is clean, modular, and maintainable.
