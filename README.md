# IDE Config Sync
A tool to save and restore IDE settings (`.vscode`, `.idea`, ...) inside Git repositories, which excludes those files.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<br>

## Setup
1. Download the [latest release](https://github.com/MatthiasHarzer/ide-config-sync/releases) and add the executable to your `PATH`.
2. Create an empty Git repository as a Database on an online service (e.g. GitHub, GitLab, Bitbucket) to store and sync your IDE settings.
3. Run `ide-config-sync init --url <URL>` with the URL of the repository from step 2.

## Usage
### Save IDE settings
1. Run `ide-config-sync save` to crawl the current directory recursively and save all IDE settings inside the Database-Repository.

### Restore IDE settings
1. Run `ide-config-sync restore` to restore all IDE settings from the Database-Repository to the current directory.

## Supported IDEs
- Visual Studio Code (`.vscode`)
- JetBrains IDEs (`.idea`)

## How it works
1. The tool searches the directory for Git repositories and their associated remotes.
2. All IDE specific folders inside the repository are saved to the Database-Repository, related to the repositories remote URL.
3. To restore the IDE settings, the tool searches the Database-Repository for the remote URL of the current directory and restores the IDE settings.