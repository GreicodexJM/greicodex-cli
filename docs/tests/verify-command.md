# E2E Test Case: `grei verify`

## Scenario: Successfully verify a compliant project

**Given** a developer is in the root directory of a project initialized with `grei init`
**And** the project contains a valid `grei.yml` file
**And** the project structure adheres to the Greicodex standards

**When** the developer runs `grei verify`

**Then** the CLI should read the `grei.yml` file
**And** check the project structure for required files and directories
**And** perform all standard checks (e.g., linting, secret scanning)
**And** display a success message indicating that the project is compliant
**And** exit with a zero status code.

## Scenario: Fail verification if `grei.yml` is missing

**Given** a developer is in a directory that does not contain a `grei.yml` file

**When** the developer runs `grei verify`

**Then** the CLI should detect that the `grei.yml` file is missing
**And** display an error message
**And** exit with a non-zero status code.

## Scenario: Fail verification if a required file is missing

**Given** a project is missing a `README.md` file

**When** the developer runs `grei verify`

**Then** the CLI should detect the missing `README.md` file
**And** display an error message listing the missing file
**And** exit with a non-zero status code.

## Scenario: Fail verification on linting errors

**Given** a project contains a Go file with linting errors

**When** the developer runs `grei verify`

**Then** the CLI should run the linter
**And** detect the linting errors
**And** display a report of the errors
**And** exit with a non-zero status code.

## Scenario: Malformed `grei.yml` file

**Given** a project's `grei.yml` file is malformed (e.g., invalid YAML syntax)

**When** the developer runs `grei verify`

**Then** the CLI should attempt to parse the `grei.yml` file
**And** fail due to the malformed content
**And** display a clear error message
**And** exit with a non-zero status code.

## Scenario: Run verification with JSON output

**Given** a compliant project

**When** the developer runs `grei verify --output=json`

**Then** the CLI should perform all verification checks
**And** output the results in a machine-readable JSON format
**And** the JSON output should indicate that all checks passed
**And** exit with a zero status code.

## Scenario: Fail verification with JSON output

**Given** a project with linting errors

**When** the developer runs `grei verify --output=json`

**Then** the CLI should perform all verification checks
**And** output the results in a machine-readable JSON format
**And** the JSON output should contain details about the linting errors
**And** exit with a non-zero status code.
