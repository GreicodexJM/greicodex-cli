# E2E Test Case: `grei install-hooks`

## Scenario: Successfully install Git hooks in a project

**Given** a developer is in the root directory of a Git repository
**And** the repository has a `.git` directory

**When** the developer runs `grei install-hooks`

**Then** the CLI should create a `pre-commit` hook in the `.git/hooks` directory
**And** the `pre-commit` hook should be executable
**And** the hook should be configured to run `grei verify`
**And** the command should exit with a success message.

## Scenario: Fail if not in a Git repository

**Given** a developer is in a directory that is not a Git repository (i.e., no `.git` directory)

**When** the developer runs `grei install-hooks`

**Then** the CLI should detect that it is not in a Git repository
**And** display an error message
**And** exit with a non-zero status code.

## Scenario: Hooks are already installed

**Given** a Git repository already has a `pre-commit` hook installed by `grei`

**When** the developer runs `grei install-hooks` again

**Then** the CLI should detect the existing hook
**And** inform the user that the hooks are already installed
**And** exit gracefully with a zero status code.

## Scenario: Overwrite existing hooks with --force flag

**Given** a Git repository has an existing `pre-commit` hook (not installed by `grei`)

**When** the developer runs `grei install-hooks --force`

**Then** the CLI should overwrite the existing `pre-commit` hook
**And** install the `grei` hook
**And** exit with a success message.

## Scenario: Fail to install hooks due to permissions error

**Given** a developer does not have write permissions for the `.git/hooks` directory

**When** the developer runs `grei install-hooks`

**Then** the CLI should attempt to create the hook file
**And** fail due to the permissions error
**And** display a clear error message
**And** exit with a non-zero status code.
