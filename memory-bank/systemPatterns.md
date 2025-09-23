# 📘 System Patterns — GRX CLI

## 1. Architecture Overview
The GRX CLI is designed using **Hexagonal Architecture (Ports and Adapters)** to ensure a clear separation between the core application logic and external dependencies like the user interface (CLI), file system, and external tools. This approach enhances modularity, testability, and maintainability.

```mermaid
graph TD
    subgraph External
        A[CLI User]
        B[File System]
        C[Plugins e.g., grei-helm]
        D[CI/CD Pipeline]
    end

    subgraph Adapters
        E[CLI Commands (Cobra)]
        F[File System Repository]
        G[Plugin Executor]
        H[JSON/ANSI Reporter]
    end

    subgraph Ports
        I[Inbound Port: CLIService]
        J[Outbound Port: FSRepository]
        K[Outbound Port: PluginRepository]
        L[Outbound Port: ReportPresenter]
    end

    subgraph Core Logic (Domain)
        M[ProjectInitializer]
        N[ProjectVerifier]
        O[HookInstaller]
        P[Scaffolder]
    end

    A -- invokes --> E
    D -- invokes --> E
    E -- calls --> I

    I -- uses --> M
    I -- uses --> N
    I -- uses --> O
    I -- uses --> P

    M -- uses --> J
    N -- uses --> J
    O -- uses --> J
    P -- uses --> J
    N -- uses --> K

    J -- implemented by --> F
    K -- implemented by --> G
    L -- implemented by --> H

    F -- interacts with --> B
    G -- interacts with --> C
    H -- outputs to --> A
    H -- outputs to --> D

    M -- notifies --> L
    N -- notifies --> L
```

## 2. Core Components (Domain)
The core contains the pure business logic of the application, with no dependencies on external technologies.
- **`ProjectInitializer`**: Handles the logic for creating a new project structure from embedded templates.
- **`ProjectVerifier`**: Contains the rules and logic for auditing a project, including secret scanning, linting, and coverage analysis.
- **`HookInstaller`**: Manages the logic for configuring Git hooks.
- **`Scaffolder`**: Implements the logic for injecting specific templates (e.g., a new stack) into an existing project.

## 3. Ports (Interfaces)
Ports are the interfaces that define the communication contracts between the core and the outside world.
- **Inbound Ports**:
  - **`CLIService`**: Defines the entry points for all user-driven actions, such as `init`, `verify`, and `scaffold`.
- **Outbound Ports**:
  - **`FSRepository`**: Defines operations for interacting with the file system, such as reading, writing, and creating files/directories.
  - **`PluginRepository`**: Defines how to discover and execute external plugins.
  - **`ReportPresenter`**: Defines the interface for presenting results to the user, either in ANSI or JSON format.

## 4. Adapters (Implementations)
Adapters are the concrete implementations of the ports that bridge the core logic with external systems.
- **Inbound Adapters**:
  - **CLI Commands**: Implemented using a library like `Cobra`, these adapters parse user input and call the `CLIService`.
- **Outbound Adapters**:
  - **`FileSystemRepository`**: The concrete implementation that uses Go's `os` and `io/fs` packages to interact with the local file system.
  - **`PluginExecutor`**: Discovers executables named `grei-<plugin>` in the system's `PATH` and communicates with them over `stdin/stdout` using a JSON protocol.
  - **`JSON/ANSI Reporter`**: Implements the `ReportPresenter` to format output for the console or for machine consumption.

## 5. Key Technical Decisions
- **Dependency Inversion**: The core logic depends on abstractions (ports), not on concrete implementations (adapters). This allows for easy testing and swapping of external dependencies. For example, the `FileSystemRepository` can be replaced with an in-memory mock for unit tests.
- **Embedded Templates**: All templates will be embedded directly into the Go binary using `go:embed`. This ensures the CLI is fully functional offline. The system will also support overriding these templates with versions from `~/.config/grei/templates` or a local `.grei/templates` directory.
- **Plugin Protocol**: Plugins are standalone executables that communicate with the main CLI via a simple JSON-based protocol on `stdin` and `stdout`. This makes the plugin system language-agnostic.
