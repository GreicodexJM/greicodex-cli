# Prompt for Generating New Greicodex CLI Templates

## Goal

Generate a new set of templates for the Greicodex CLI, based on a specified technology stack. The templates should be consistent with the project's standards for hexagonal architecture, local development, testing, linting, and deployment.

## Context

The Greicodex CLI is a command-line tool that automates the initialization and standardization of software projects. It uses a set of embedded templates to scaffold new projects based on a user-defined recipe. The templates are written in Go's `text/template` format, with placeholders for project-specific variables (e.g., `{{ .Project.Name }}`).

## Instructions

1.  **Create the Template Directory**: Create a new directory for the template at `internal/core/scaffolder/templates/{language}-{tooling}-{...other-properties}`. For example, a template for a React project with an Express NodeJS backend and postgresql persistence would be located at `internal/core/scaffolder/templates/typescript-react-express-pgsql`.

2.  **Populate the Template Files**: Create the following files in the new template directory, using the provided content as a guide.

    *   **`package.json.tmpl`**: Define the project's dependencies, including the core framework, TypeScript, ESLint, and Jest.
    *   **`tsconfig.json.tmpl`**: Configure the TypeScript compiler options.
    *   **`.eslintrc.js.tmpl`**: Define the ESLint rules for the project.
    *   **`jest.config.js.tmpl`**: Configure the Jest testing framework.
    *   **`Dockerfile.tmpl`**: Create a multistage Dockerfile that builds the application and serves it with Nginx. Include a debug stage with the appropriate tools (e.g., Delve for Go, Node Inspector for TypeScript).
    *   **`docker-compose.yml.tmpl`**: Define the local development environment, including an `app` service for running the application and a `debug` service for debugging.
    *   **`.vscode/launch.json.tmpl`**: Configure the debugger for Visual Studio Code.
    *   **Hexagonal Architecture Directories**: Create the following directories and add a `.gitkeep` file to each:
        *   `src/domain`
        *   `src/ports/inbound`
        *   `src/ports/outbound`
        *   `src/adapters/inbound`
        *   `src/adapters/outbound`
    *   **Basic Application Structure**: Create the basic application structure, including a main entry point and a sample component.

3.  **Create a `manifest.yml` File**: Create a `manifest.yml` file in the root of the new template directory. This file should contain all the metadata for the stack, including its name, description, type, and the technologies it provides.

## Example

Here's an example of the `Dockerfile.tmpl` for the `go-cobra` template:

```Dockerfile
# Build stage
FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /{{ .Project.Name }} ./cmd/app

# Debug stage
FROM golang:1.18-alpine AS debug

WORKDIR /app

COPY --from=build /{{ .Project.Name }} /{{ .Project.Name }}

RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 40000

CMD ["dlv", "--listen=:40000", "--headless", "--api-version=2", "exec", "/{{ .Project.Name }}"]

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /{{ .Project.Name }} /{{ .Project.Name }}

USER nonroot:nonroot

ENTRYPOINT ["/{{ .Project.Name }}"]
