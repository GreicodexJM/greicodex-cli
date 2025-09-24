# 📘 Progress — GRX CLI

## 1. Current Status
**Phase 1 - Core CLI (In Progress)**

The plugin system is now fully functional and includes built-in support for common tech stacks. The `init` command can now offer Symfony/LAMP and MERN stacks out of the box, in addition to any external plugins the user has installed. The survey logic correctly de-duplicates project types from multiple plugins, and the generated `grei.yml` is now richly populated from the chosen plugin's manifest.

### What Works
- **Project Documentation:** The `memory-bank` is fully initialized and up-to-date.
- **Plugin Discovery:**
  - [x] A `PluginScanner` finds both built-in and external plugins.
- **`grei init` command:**
  - [x] Dynamically builds the survey from all available plugins.
  - [x] De-duplicates project types.
  - [x] Fully populates the `grei.yml` recipe from the chosen plugin.
  - [x] Includes a "Custom" fallback.

### What's Left to Build (Phase 1)
- **Create a Reference Plugin:**
  - [ ] Build a simple `grei-mock-symfony` external plugin to test the system.
- **`grei verify` command:**
  - [ ] Read the `grei.yml` file to become context-aware.
- **Template Scaffolding:**
  - [ ] Implement logic for the core CLI to invoke a `scaffold` command on the chosen plugin.

## 2. Known Issues
- None at this time.

## 3. Evolution of Project Decisions
- **Pivotal Change:** The CLI now supports a **hybrid plugin model**.
  - **Reasoning:** Combining built-in plugins with an extensible external plugin system provides the best of both worlds: immediate utility and long-term flexibility.
