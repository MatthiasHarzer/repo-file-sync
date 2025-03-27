# Repo File Sync
A tool to save and restore ignored files inside a Git repositories, based on defined Glob-patterns.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<br>

## Setup 
Download the [latest release](https://github.com/MatthiasHarzer/repo-file-sync/releases) and add the executable to your `PATH`.

### Local only
This will save the files local only, without the need of a remote repository.
1. Run `repo-file-sync init`.
2. Optionally provide a directory on the machine to store the Database-Repository, however the default one is recommended.
3. Select `y` when asked to set up local only
4. Follow the instructions to finish the setup.

### Remote repository
This will save the files to a remote repository to sync them between multiple devices.
1. Run `repo-file-sync init` and provide a directory to store the Database-Repository.
2. Optionally provide a directory on the machine to store the Database-Repository, however the default one is recommended.
3. Select `n` when asked to set up local only
4. Provide the URL of the remote repository to store the IDE configurations
    - For that, set up a new repository on GitHub, GitLab, Bitbucket, etc.
    - Copy the URL of the repository and paste it into the prompt
5. Follow the instructions to finish the setup.


## Usage
> All following commands support the `--dir <directory>` / `-d <directory>` flag to search and save/restore files in a specific directory. If no directory is specified, the current directory is used.

### Save files
1. Run `repo-file-sync save` to crawl the directory recursively and save all matching files inside the Database-Repository.

### Restore files
1. Run `repo-file-sync restore` to search the directory recursively and restore all files from the Database-Repository.

### Dry run / discover
1. Run `repo-file-sync discover` to see which files would be saved or restored.

## Included / excluded files
The tool uses a set of global include/exclude patterns to determine which files to save and restore. Additionally, every repository can have its own include/exclude patterns.
### Included files
Repo File Sync was originally designed to save and restore IDE configurations, so it includes by default the following patterns:
- Visual Studio Code (`**/.vscode/**`)
- JetBrains IDEs (`**/.idea/**`)

#### Adding include patterns
1. Run `repo-file-sync pattern include add <pattern-1>, <pattern-2>, ...` to add one or more Glob-patterns to match files from the repository root.

#### Removing include patterns
1. Run `repo-file-sync pattern include remove <pattern-1>, <pattern-2>, ...` to remove one or more Glob-patterns.

> To add or remove a global include pattern, use the `--global/-g` flag.

### Excluded files
The following patterns are excluded by default, to prevent matching files in some ignored folders:
- Node modules (`**/node_modules/**`)
- Virtual environments (`**/venv/**`, `**/.venv/**`)
> Exclude patterns are evaluated **after** include patterns, so a file matching an exclude pattern will always be excluded, even if it matches an include pattern.

#### Adding exclude patterns
1. Run `repo-file-sync pattern exclude add <pattern-1>, <pattern-2>, ...` to add one or more Glob-patterns to exclude files from the repository root.

#### Removing exclude patterns
1. Run `repo-file-sync pattern exclude remove <pattern-1>, <pattern-2>, ...` to remove one or more Glob-patterns.

> To add or remove a global exclude pattern, use the `--global/-g` flag.


### Excluded folders when searching for repositories
The tool excludes the following folders when discovering repositories:
- Node modules (`node_modules`)
- Virtual environments (`venv`, `.venv`)
> This has no impact on the saved and restored files. Use the include and exclude patterns to define which files to save and restore.