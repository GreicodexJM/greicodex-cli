# 📘 Tech Context — GRX CLI

## 1. Core Technology Stack
- **Language:** **Go (Golang)**
- **CLI Framework:** **`Cobra`**
- **Template Management:** **`go:embed`**

## 2. Development and Build Setup
- **Project Structure:** Standard Go project layout.
- **Dependency Management:** Go Modules.
- **Distribution:** **`GoReleaser`**.

## 3. The Internal Stack Registry
The GRX CLI features a **compositional internal stack architecture**. All available technology stacks are defined and managed directly within the core CLI binary in a central registry.

### Stack Definition
Each built-in stack is defined as a `stack.Stack` struct in `internal/core/stack/registry.go`. This struct contains:
- **Name:** A unique identifier for the stack (e.g., `symfony-lamp`).
- **Description:** A user-facing description.
- **Type:** The category of the stack (`code`, `persistence`, `deployment`).
- **Provides:** A struct detailing the specific technologies (Language, Tooling, etc.) that this stack provides.

### Extensibility Model
To add a new technology stack to the CLI, a developer must:
1.  Add a new `stack.Stack` definition to the `internal/core/stack/registry.go` file.
2.  Implement the necessary logic for template scaffolding and verification for the new stack.
3.  Submit a pull request to the core `grei-cli` repository.

This model ensures that all available stacks are centrally managed, versioned, and vetted according to Greicodex standards, maintaining the goal of a consistent and high-quality single-binary solution.

## 4. Tool Usage Patterns
- **Template Overrides:** The CLI supports overriding embedded templates via project-specific or user-specific directories.
- **Docker Compose Handling:** The CLI respects the `docker-compose.yml` and `docker-compose.override.yml` convention.
