# 📘 Active Context — GRX CLI

## 1. Current Work Focus
The plugin system is now robust, with the `init` command capable of fully populating the `grei.yml` recipe from a chosen plugin's manifest. The next logical step is to create a reference plugin to test the end-to-end flow.

**Immediate Goals:**
1.  **Create a Mock Plugin:** Develop a simple mock external plugin (e.g., `grei-mock-symfony`) to test the discovery and scaffolding process end-to-end.
2.  **Integrate `grei.yml` into the `verify` command.**
3.  **Implement Template Scaffolding:** Move beyond the placeholder and implement the actual logic for copying templates based on the recipe.

## 2. Recent Changes & Decisions
- **Decision:** The `init` command now populates the entire `recipe.Stack` and relevant sub-structs (`WebApp`, `Api`) from the chosen plugin manifest.
- **Reasoning:** This makes the generated `grei.yml` a much richer and more accurate representation of the intended project architecture.
- **Decision:** Added **built-in plugins** for Symfony/LAMP and MERN stacks.
- **Decision:** The `init` command is now fully **plugin-driven**.

## 3. Next Steps
- **Begin work on a simple external plugin** to serve as a reference implementation.
- **Refactor `verifier.Service`** to be recipe-aware.
- **Implement the `scaffold` command** on plugins.

## 4. Important Patterns & Preferences
- **Hybrid Plugin Model:** The CLI should support both built-in and external plugins.
- **Standardized Discovery:** All plugins must adhere to the `discover` command protocol.
