# 📘 Active Context — GRX CLI

## 1. Current Work Focus
The primary focus is on **Phase 1: Core CLI Implementation**. With the `init` command now reworked, the next major goal is to make the `verify` command context-aware by reading the `grei.yml` recipe file.

**Immediate Goals:**
1.  **Integrate `grei.yml` into the `verify` command:** The verifier service should read the recipe to determine which checks to run (e.g., which linter config to look for).
2.  **Implement Linter Detection:** Add the logic to the `verify` command to check for the existence of the linter configuration file specified in `grei.yml`.
3.  **Expand Template Scaffolding:** Move beyond the placeholder and implement the actual logic for copying templates based on the recipe.

## 2. Recent Changes & Decisions
- **Decision:** Implemented the **Project Recipe** concept. The `grei init` command is now an interactive TUI that generates a `grei.yml` file. This file acts as the single source of truth for a project's configuration.
- **Reasoning:** This was a critical change to ensure that the CLI's behavior is driven by explicit, version-controlled configuration, making it more robust and intelligent. It also provides essential context for AI agents.
- **Dependencies Added:**
  - `github.com/AlecAivazis/survey/v2` for the interactive TUI.
  - `github.com/briandowns/spinner` for progress indicators.
  - `gopkg.in/yaml.v3` for YAML marshalling.

## 3. Next Steps
- **Refactor `verifier.Service`:** Modify the service to accept the `recipe.Recipe` as an input.
- **Create Linter Port/Adapter:** Design and implement the `LinterDetector` port and its corresponding filesystem adapter.
- **Update `verify` command:** Wire the new components into the `verify` command adapter.

## 4. Important Patterns & Preferences
- **Recipe-Driven:** All commands should, where applicable, consult the `grei.yml` file to inform their behavior.
- **Spanish-first:** All user-facing output should be in Spanish.
- **Offline-first:** The CLI must be fully functional without an internet connection.
