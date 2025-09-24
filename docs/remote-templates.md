# Remote Templates

The GRX CLI fetches project templates from a remote Git repository. This ensures that you always have access to the latest templates without needing to update the CLI tool itself.

## How it Works

When you run the `grei init` command, the CLI performs the following steps:

1.  **Clones or Pulls the Template Repository**: The CLI clones or pulls the latest version of the [greicodex-cli repository](https://github.com/GreicodexJM/greicodex-cli.git) into a local cache directory. This ensures that you have the most up-to-date templates.

2.  **Version Check**: The CLI reads a `manifest.json` file from the `templates` directory of the repository. This file specifies the minimum required version of the GRX CLI to use the templates. If your CLI version is older than the required version, the `init` command will fail with an error message.

3.  **Project Initialization**: If the version check passes, the CLI proceeds to initialize your project using the downloaded templates.

## Offline Support

The CLI caches the downloaded templates locally. If you run `grei init` without an internet connection, it will use the cached templates. This allows you to initialize projects even when you are offline, as long as you have run the command at least once with an internet connection.
