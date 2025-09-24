# AmableSagitta

This project was initialized by the GRX CLI.


# 📘 Project Documentation — GRX CLI

---

## 1. Project Brief

**Project name:** GRX CLI (Greicodex CLI)
**Goal:** Build a Spanish-first command-line tool that automates initialization, verification, and standardization of software projects under Greicodex guidelines. It should allow any developer to set up a new project or audit an existing one easily, reproducibly, and offline-friendly.
**Motivation:**

* Reduce variability in project configuration.
* Enforce good practices (Hexagonal architecture, SOLID, 12-Factor principles, CI/CD).
* Avoid common mistakes in local vs QA/PROD environments.
* Integrate with existing tooling (e.g., Docker Compose → Helm converter).

---

## 2. Project Scope

### In Scope

* Lightweight CLI binary with commands in Spanish.
* Initialization of repositories with standardized templates (README, hooks, pipelines, docker-compose, Helm, IaC).
* Verification of existing repositories (secret scan, lint, tests, coverage, Helm, IaC).
* Automatic installation of Git hooks.
* Scaffolding of CI/CD, Helm charts, and OpenTofu templates.
* Modular support through **plugins** (`grei-<plugin>`).
* ANSI/JSON output for human-friendly and CI integrations.
* Handling of `docker-compose.yml` (base) and `docker-compose.override.yml` (local), respecting safe defaults.
* Ability to work **offline** (embedded templates).

### Out of Scope

* No custom CI/CD system will be built, only scaffolding.
* No user authentication or role management.
* The AI agent is not part of the CLI; it is the one building the CLI.

---

## 3. Functional Requirements

### Core Commands

1. `grei init [path]`

   * Creates standardized project structure.
   * Copies templates (README, pipelines, docker-compose, Helm, IaC).
   * Prepares `.githooks` and configures Git.

2. `grei verify [path] [--min-cov=N]`

   * Runs secret scanning (gitleaks if present, regex fallback).
   * Runs linters if config is detected.
   * Runs tests and calculates minimum coverage.
   * Verifies CI/CD, Helm, and OpenTofu configuration.
   * Reports in ANSI (human) or JSON (`--json`).

3. `grei install-hooks [path]`

   * Configures `core.hooksPath`.
   * Marks hooks as executable.
   * Tests hook activation.

4. `grei scaffold <ci|helm|tofu|stack>`

   * Injects templates on demand.
   * Example: `grei scaffold stack ts` generates TypeScript setup (Jest + ESLint).

5. `grei plugin <list|run>`

   * Lists detected plugins in `PATH`.
   * Executes plugin via JSON protocol (`stdin/stdout`).

---

## 4. Tech Stack

* **Core language:** Go (Golang).
* **Distribution:** single binary (`goreleaser`), `.deb`, Homebrew tap, tar.gz.
* **Templates:** embedded via `go:embed`, with overrides in `~/.config/grei/templates` and `.grei/templates`.
* **Plugins:** `grei-<plugin>` executables discovered in `PATH`.
* **Coverage parsing:** JSON (`coverage-summary.json`), XML (Cobertura/JaCoCo/phpunit/pytest).
* **ANSI output:** lightweight library (`fatih/color`) or manual escape codes.
* **Secret scan:** integration with `gitleaks` + regex fallback.
* **Future integrations:**

  * Python-based tool for converting `docker-compose` → Helm (as plugin).
  * Handling of `docker-compose.yml` + `docker-compose.override.yml` by default.

---

## 5. Execution Plan

### Phase 1 — Core CLI (2 weeks)

* Implement `init`, `verify`, `install-hooks`.
* Embed templates.
* Coverage parser.
* ANSI/JSON output.

### Phase 2 — Plugins & Scaffolding (2–3 weeks)

* Add plugin discovery.
* Stack scaffolding (TS, PHP, Python, Go, Java, Dart).
* Integrate Python `compose → helm` converter as a plugin.

### Phase 3 — Packaging & Distribution (1 week)

* `goreleaser` for deb/rpm/tar.gz/homebrew.
* Installation script (`curl | sh`).
* Usage documentation.

### Phase 4 — Extensions (optional)

* `self-update`.
* JIRA/Cronus integrations (reminders).
* Offline caches (`npm`, `pip`, `composer`).

---

## 6. Current Scope Context

* **Existing tools:**

  * Python utility already exists to convert `docker-compose.yml` to Helm.
  * Internal practice:

    * `docker-compose.yml` describes **base services** (QA/PROD).
    * `docker-compose.override.yml` describes **local development environment**.
    * This avoids accidents by omission and enforces safe defaults.

* **Integration strategy:**

  * GRX CLI core will respect this convention (`compose` base + override).
  * The Helm conversion plugin will use **docker-compose.yml** as primary input, with optional overlays.

