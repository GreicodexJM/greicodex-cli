# 📘 Product Context — GRX CLI

## 1. Problem Statement
Developers at Greicodex face inconsistencies when setting up new projects or auditing existing ones. This variability leads to configuration drift, repeated mistakes, and a lack of standardization, especially between local, QA, and production environments. Key challenges include:
- **Manual setup is error-prone:** Manually configuring linters, Git hooks, CI/CD pipelines, and infrastructure-as-code (IaC) for each project is time-consuming and inconsistent.
- **Good practices are hard to enforce:** Ensuring every project adheres to Hexagonal Architecture, SOLID principles, and 12-Factor App guidelines requires constant oversight.
- **Environment discrepancies:** Differences between `docker-compose.yml` (base services) and `docker-compose.override.yml` (local development) often cause "works on my machine" issues.

## 2. How It Should Work
The GRX CLI will be a Spanish-first, offline-friendly command-line tool that acts as a "project guardian." It will provide a simple, unified interface to:
- **Initialize projects:** `grei init` will create a standardized repository structure with pre-configured templates for READMEs, Git hooks, CI/CD pipelines, Docker Compose, Helm charts, and OpenTofu.
- **Verify compliance:** `grei verify` will audit a project against Greicodex standards, checking for secrets, code quality, test coverage, and valid IaC configurations.
- **Scaffold components:** `grei scaffold` will inject specific templates on demand, such as a TypeScript stack with Jest and ESLint, or a Helm chart for a new service.
- **Extend functionality:** A modular plugin system (`grei-<plugin>`) will allow for easy integration of additional tools, like the existing Python-based `docker-compose` to Helm converter.

## 3. User Experience Goals
- **Simplicity:** The CLI should be intuitive, with clear commands and helpful feedback in Spanish.
- **Reproducibility:** Any developer should be able to achieve the same project setup by running the same commands.
- **Offline-first:** The tool must be fully functional without an internet connection, using embedded templates.
- **CI/CD friendly:** Output should be available in both human-readable ANSI format and machine-readable JSON format for easy integration into automated pipelines.
