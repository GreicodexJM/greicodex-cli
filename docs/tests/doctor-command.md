# E2E Test Case: `grei doctor`

## Scenario: All system checks pass

**Given** a developer has all the required system dependencies installed (e.g., Git, Docker)
**And** all dependencies are at a compatible version

**When** the developer runs `grei doctor`

**Then** the CLI should check for the presence and version of each required dependency
**And** display a success message for each check
**And** conclude with a summary message indicating that the system is ready
**And** exit with a zero status code.

## Scenario: A required dependency is missing

**Given** a developer does not have Docker installed

**When** the developer runs `grei doctor`

**Then** the CLI should check for all required dependencies
**And** detect that Docker is not installed
**And** display an error message for the missing dependency
**And** display success messages for the dependencies that are correctly installed
**And** conclude with a summary message indicating that the system is not ready
**And** exit with a non-zero status code.

## Scenario: An installed dependency is outdated

**Given** a developer has an outdated version of Git installed

**When** the developer runs `grei doctor`

**Then** the CLI should check for all required dependencies
**And** detect that the Git version is below the minimum requirement
**And** display an error message indicating the version incompatibility
**And** conclude with a summary message indicating that the system is not ready
**And** exit with a non-zero status code.

## Scenario: Run doctor with JSON output

**Given** a developer has some dependencies installed correctly and some missing

**When** the developer runs `grei doctor --output=json`

**Then** the CLI should perform all system checks
**And** output the results in a machine-readable JSON format
**And** the JSON output should clearly indicate the status of each dependency
**And** exit with a non-zero status code.

## Scenario: Command not found for a dependency check

**Given** a dependency check relies on a command that is not in the system's PATH

**When** the developer runs `grei doctor`

**Then** the CLI should handle the "command not found" error gracefully
**And** report that the specific dependency check failed
**And** continue with the other checks
**And** exit with a non-zero status code.
