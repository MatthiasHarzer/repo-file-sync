# IDE Config Sync
A tool to save and restore IDE settings (`.vscode`, `.idea`) inside Git repositories, which excludes those files.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<br>

## Setup 
Download the [latest release](https://github.com/MatthiasHarzer/ide-config-sync/releases) and add the executable to your `PATH`.

### Local only
This will save the IDE configurations local only, without the need of a remote repository.
1. Run `ide-config-sync init`.
2. Optionally provide a directory on the machine to store the Database-Repository, however the default one is recommended.
3. Select `y` when asked to set up local only

### Remote repository
This will save the IDE configurations to a remote repository to sync them between multiple devices.
1. Run `ide-config-sync init` and provide a directory to store the Database-Repository.
2. Optionally provide a directory on the machine to store the Database-Repository, however the default one is recommended.
3. Select `n` when asked to set up local only
4. Provide the URL of the remote repository to store the IDE configurations
    - For that, set up a new repository on GitHub, GitLab, Bitbucket, etc.
    - Copy the URL of the repository and paste it into the prompt

## Usage
> All following commands support the `--dir <directory>` / `-d <directory>` flag to search and save/restore IDE settings in a specific directory. If no directory is specified, the current directory is used.

### Save IDE settings
1. Run `ide-config-sync save` to crawl the directory recursively and save all IDE settings inside the Database-Repository.

### Restore IDE settings
1. Run `ide-config-sync restore` to search the directory recursively and restore the IDE settings from the Database-Repository.

### Dry run / discover
1. Run `ide-config-sync discover` to see which IDE settings would be saved or restored.

## Supported IDEs
- Visual Studio Code (`.vscode`)
- JetBrains IDEs (`.idea`)

## How it works
1. The tool searches the directory for Git repositories and their associated remotes.
2. All IDE specific folders inside the repository are saved to the Database-Repository, related to the repositories remote URL.
3. To restore the IDE settings, the tool searches the Database-Repository for the remote URL of the current directory and restores the IDE settings.