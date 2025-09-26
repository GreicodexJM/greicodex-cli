# Prompt for Generating New Greicodex CLI Templates

## Goal

Generate a new set of templates for the Greicodex CLI, based on a specified technology stack. The templates should be consistent with the project's standards for hexagonal architecture, local development, testing, linting, and deployment.

## Context

The Greicodex CLI is a command-line tool that automates the initialization and standardization of software projects. It uses a set of embedded templates to scaffold new projects based on a user-defined recipe. The templates are written in Go's `text/template` format, with placeholders for project-specific variables (e.g., `{{ .Project.Name }}`).

## Instructions

1.  **Create the Template Directory**: Create a new directory for the template at `templates/skeletons/{stack-name}`. For example, a template for a Typescript Express API would be located at `templates/skeletons/typescript-express`.

2.  **Create a `manifest.yml` File**: Create a `manifest.yml` file in the root of the new template directory. This file should contain all the metadata for the stack, including its name, description, type, the technologies it provides, and the options that will be presented to the user in the TUI.

    ```yaml
    name: typescript-express
    description: Pila para crear una API en Typescript con Express.
    type: api
    provides:
      language: Typescript
      tooling: Express
    options:
      persistence:
        message: "¿Qué tipo de persistencia usarás?"
        values:
          - "None"
          - "PostgreSQL"
          - "MySQL"
      deployment:
        message: "¿Dónde quieres desplegar la aplicación?"
        values:
          - "Kubernetes"
          - "Lambda"
    ```

3.  **Populate the Template Files**: Create the template files in the new directory. Use the `.tmpl` extension for any file that needs to be processed by the template engine. You can use placeholders and expressions to make the templates dynamic.

    For example, you can use `{{ .Project.Name }}` to insert the project name, and `{{ if eq .Stack.deployment "Kubernetes" }}` to conditionally include content based on the user's selections.

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
