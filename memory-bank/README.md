# 📘 Project Brief — GRX CLI

## 1. Project Name
GRX CLI (Greicodex CLI)

## 2. Goal
Build a Spanish-first command-line tool that automates the initialization, verification, and standardization of software projects under Greicodex guidelines. It should allow any developer to set up a new project or audit an existing one easily, reproducibly, and in an offline-friendly manner.

## 3. Motivation
- **Reduce variability:** Minimize inconsistencies in project configuration across different teams and developers.
- **Enforce good practices:** Systematically apply principles like Hexagonal Architecture, SOLID, and 12-Factor App methodology, along with CI/CD best practices.
- **Prevent common errors:** Avoid typical mistakes that arise from discrepancies between local, QA, and production environments.
- **Integrate with existing tooling:** Leverage and connect with current tools, such as a Docker Compose to Helm converter, to streamline workflows.

## 4. New Features
### SPA Templates
The CLI now supports scaffolding for Single Page Applications (SPAs) using Angular and Vue.js. These templates come pre-configured with:
- TypeScript support
- ESLint for code linting
- Jest for unit testing and test coverage
- Multistage Dockerfiles for building and serving with Nginx
- CI/CD deployment options for Netlify and Kubernetes with Knative Serving
- Hooks for secret scanning and other checks
- Hexagonal architecture structure for organizing code
