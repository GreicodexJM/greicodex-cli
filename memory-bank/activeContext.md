# 📘 Active Context — GRX CLI

## 1. Current Work Focus
Test coverage has been added to the project's `Makefile`, improving the project's own quality and standards. The immediate focus is to continue expanding the `verify` command's checks.

**Immediate Goals:**
1.  **Expand the `verify` command's checks** further (e.g., deployment checks).
2.  **Add more built-in stacks** and templates to the registry.

## 2. Recent Changes & Decisions
- **Decision:** Added **test coverage** generation to the project's `Makefile`.
- **Reasoning:** This is a critical step in "dogfooding" the CLI and ensuring it adheres to the same high standards it will enforce on other projects.
- **Decision:** Expanded the `verify` command to include checks for persistence and linters.
- **Decision:** Removed the external plugin system in favor of a **Compositional Internal Stack Architecture**.

## 3. Next Steps
- **Add a deployment check** to the `verifier.Service` that is driven by the `recipe.Deployment.Type` field.
- **Add a `deployment` stack** for Serverless to the internal registry.

## 4. Important Patterns & Preferences
- **Dogfooding:** The CLI should be capable of managing its own project structure and standards.
- **Single Binary:** The CLI should be a self-contained, single executable.
