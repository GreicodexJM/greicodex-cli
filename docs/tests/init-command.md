# E2E Test Case: `grei init`

## Scenario: Successfully initialize a new project

**Given** a developer is in an empty directory
**And** they have a stable internet connection
**And** the GRX CLI version is compatible with the remote templates

**When** the developer runs `grei init my-new-project`
**And** follows the interactive prompts to select a "go-cli" stack
**And** selects "Ninguna" for persistence
**And** selects "Ninguno" for deployment

**Then** the CLI should create a new directory named `my-new-project`
**And** inside `my-new-project`, there should be a `grei.yml` file with the selected configuration
**And** the directory should contain the standard project structure, including:
  - `README.md`
  - `.gitignore`
  - `LICENSE`
  - `docs/`
  - `deploy/helm/`
**And** the CLI should initialize a new Git repository in the `my-new-project` directory
**And** a `develop` branch should be created
**And** the command should exit with a success message.

## Scenario: Initialize a project in the current directory

**Given** a developer is in an empty directory `my-current-project`
**And** they have a stable internet connection

**When** the developer runs `grei init .`
**And** completes the interactive survey

**Then** the CLI should use the current directory (`my-current-project`) as the project root
**And** all standard project files and directories should be created in the current directory
**And** the command should exit with a success message.

## Scenario: Fail initialization if project already exists

**Given** a directory `my-existing-project` already contains a `grei.yml` file

**When** the developer runs `grei init my-existing-project`

**Then** the CLI should detect the existing `grei.yml` file
**And** display an error message indicating that a project already exists
**And** exit with a non-zero status code without modifying any files.

## Scenario: Fail initialization if CLI version is outdated

**Given** the remote template repository's `manifest.json` requires a minimum CLI version of `1.0.0`
**And** the developer is using GRX CLI version `0.1.0`

**When** the developer runs `grei init my-project`

**Then** the CLI should download the templates
**And** read the `manifest.json` file
**And** compare the versions
**And** display an error message about the outdated CLI version
**And** exit with a non-zero status code.

## Scenario: Fail initialization due to network error

**Given** a developer does not have an internet connection

**When** the developer runs `grei init my-project`

**Then** the CLI should attempt to download the remote templates
**And** fail due to the network error
**And** display a clear error message indicating the connection failure
**And** exit with a non-zero status code.

## Scenario: User cancels the interactive survey

**Given** a developer is running `grei init`

**When** the developer cancels the survey midway (e.g., by pressing `Ctrl+C`)

**Then** the CLI should exit gracefully
**And** should not create any files or directories.

## Scenario: Invalid `manifest.json` in remote repository

**Given** the `manifest.json` file in the remote repository is malformed or contains an invalid version string

**When** the developer runs `grei init my-project`

**Then** the CLI should download the templates
**And** attempt to parse the `manifest.json` file
**And** fail due to the invalid format
**And** display a clear error message
**And** exit with a non-zero status code.
