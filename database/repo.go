package database

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"repo-file-sync/repository"
	"repo-file-sync/set"
	"repo-file-sync/util/fsutil"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
)

func remoteURLToDir(remoteURL string) string {
	invalidChars := regexp.MustCompile(`[:<>"\\|?*]+`)

	remoteURL = strings.TrimPrefix(remoteURL, "https://")
	remoteURL = strings.TrimPrefix(remoteURL, "http://")
	sanitized := invalidChars.ReplaceAllString(remoteURL, "_")
	sanitized = strings.ReplaceAll(sanitized, "__", "_")
	sanitized = strings.Trim(sanitized, "/")

	return sanitized
}

type Repo struct {
	Directory            string
	repo                 *git.Repository
	changesSinceLastPush map[string]int
}

func NewRepoDatabase(directory string) (Database, error) {
	_, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Directory:            directory,
		changesSinceLastPush: make(map[string]int),
	}, nil
}

func InitializeRepoDatabaseFromURL(url, directory string) (Database, error) {
	cmd := exec.Command("git", "clone", url, directory)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return NewRepoDatabase(directory)
}

func InitializeRepoDatabaseFromPath(directory string) (Database, error) {
	_, err := git.PlainInit(directory, false)
	if err != nil {
		return nil, err
	}

	return NewRepoDatabase(directory)
}

func (d *Repo) remoteFilesDir(remote string) string {
	remoteAsPath := remoteURLToDir(remote)
	return fmt.Sprintf("%s/%s/files", d.Directory, remoteAsPath)
}

func (d *Repo) remoteIncludesFile(remote string) string {
	remoteAsPath := remoteURLToDir(remote)
	return fmt.Sprintf("%s/%s/includes", d.Directory, remoteAsPath)
}

func (d *Repo) globalIncludesFile() string {
	return fmt.Sprintf("%s/.global/includes", d.Directory)
}

func (d *Repo) globalExcludesFile() string {
	return fmt.Sprintf("%s/.global/excludes", d.Directory)
}

func (d *Repo) writeRemoteRepoFile(remote string, file repository.File) error {
	dbFilePath := fmt.Sprintf("%s/%s", d.remoteFilesDir(remote), file.PathFromRepoRoot)
	relativePath, err := filepath.Rel(d.Directory, dbFilePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path for %s: %s", dbFilePath, err)
	}

	err = copy.Copy(file.AbsolutePath, dbFilePath)
	if err != nil {
		return fmt.Errorf("failed to copy %s to %s: %s", file.AbsolutePath, dbFilePath, err)
	}

	cmd := exec.Command("git", "add", relativePath, "--force")
	cmd.Dir = d.Directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add %s to git: %s", relativePath, err)
	}

	if _, ok := d.changesSinceLastPush[remote]; !ok {
		d.changesSinceLastPush[remote] = 0
	}
	d.changesSinceLastPush[remote]++

	return nil
}

func (d *Repo) writeRepoIncludes(remote string, includes []string) error {
	includesFile := d.remoteIncludesFile(remote)

	err := fsutil.WriteFileLines(includesFile, includes)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "add", includesFile, "--force")
	cmd.Dir = d.Directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add %s to git: %s", includesFile, err)
	}

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update includes for '%s'", remote))
	cmd.Dir = d.Directory
	_ = cmd.Run()

	return nil
}

func (d *Repo) WriteRepoFile(repo string, file repository.File) error {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		err := d.writeRemoteRepoFile(remote, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Repo) WriteRepoDiscoveryOptions(repo string, options repository.DiscoveryOptions) error {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		err := d.writeRepoIncludes(remote, options.IncludePatterns.Slice())
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Repo) getComputedRepoDiscoveryOptions(repo string) (repository.DiscoveryOptions, error) {
	globalOptions, err := d.ReadGlobalDiscoveryOptions()
	if err != nil {
		return repository.DiscoveryOptions{}, err
	}

	repoOptions, err := d.ReadRepoDiscoveryOptions(repo)
	if err != nil {
		return repository.DiscoveryOptions{}, err
	}

	return globalOptions.Merge(repoOptions), nil
}

func (d *Repo) ReadRepoFiles(repo string) (<-chan repository.File, error) {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return nil, err
	}

	options, err := d.getComputedRepoDiscoveryOptions(repo)
	if err != nil {
		return nil, err
	}

	files := make(chan repository.File)

	go func() {
		defer close(files)

		for _, remote := range remotes {
			filesPath := d.remoteFilesDir(remote)
			repoFiles := repository.DiscoverRepositoryFiles(filesPath, options)
			for file := range repoFiles {
				files <- file
			}
		}
	}()

	return files, nil
}

func (d *Repo) ReadRepoDiscoveryOptions(repo string) (repository.DiscoveryOptions, error) {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return repository.DiscoveryOptions{}, err
	}

	var includes []string
	for _, remote := range remotes {
		includesPath := d.remoteIncludesFile(remote)
		exists, _ := fsutil.Exists(includesPath)
		if !exists {
			continue
		}

		lines, err := fsutil.ReadFileLines(includesPath)
		if err != nil {
			return repository.DiscoveryOptions{}, err
		}

		includes = append(includes, lines...)
	}

	return repository.NewDiscoveryOptions(includes), nil
}

func (d *Repo) readGlobalIncludes() ([]string, error) {
	exists, _ := fsutil.Exists(d.globalIncludesFile())
	if !exists {
		return nil, nil
	}

	return fsutil.ReadFileLines(d.globalIncludesFile())
}

func (d *Repo) readGlobalExcludes() ([]string, error) {
	exists, _ := fsutil.Exists(d.globalExcludesFile())
	if !exists {
		return nil, nil
	}

	return fsutil.ReadFileLines(d.globalExcludesFile())
}

func (d *Repo) globalOptionsExists() bool {
	includesExists, _ := fsutil.Exists(d.globalIncludesFile())
	excludesExists, _ := fsutil.Exists(d.globalExcludesFile())

	return includesExists && excludesExists
}

func (d *Repo) ReadGlobalDiscoveryOptions() (repository.DiscoveryOptions, error) {
	if !d.globalOptionsExists() {
		options := repository.DefaultGlobalDiscoveryOptions()
		err := d.WriteGlobalDiscoveryOptions(options)
		if err != nil {
			return repository.DiscoveryOptions{}, err
		}
		return options, nil
	}
	includeLines, inclErr := d.readGlobalIncludes()

	if inclErr != nil {
		return repository.DiscoveryOptions{}, fmt.Errorf("failed to read global includes")
	}

	options := repository.DefaultGlobalDiscoveryOptions()
	options.IncludePatterns = set.FromSlice(includeLines)

	return options, nil
}

func (d *Repo) WriteGlobalDiscoveryOptions(config repository.DiscoveryOptions) error {
	err := fsutil.WriteFileLines(d.globalIncludesFile(), config.IncludePatterns.Slice())
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "add", d.globalIncludesFile(), "--force")
	cmd.Dir = d.Directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add %s to git: %s", d.globalIncludesFile(), err)
	}

	cmd = exec.Command("git", "add", d.globalExcludesFile(), "--force")
	cmd.Dir = d.Directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add %s to git: %s", d.globalExcludesFile(), err)
	}

	cmd = exec.Command("git", "commit", "-m", "Update global includes")
	cmd.Dir = d.Directory
	_ = cmd.Run()

	return nil
}

func (d *Repo) Push() error {
	if len(d.changesSinceLastPush) == 0 {
		return nil
	}

	totalFiles := 0
	for _, count := range d.changesSinceLastPush {
		totalFiles += count
	}

	title := "Update upto "
	if totalFiles == 1 {
		title += "1 file"
	} else {
		title += fmt.Sprintf("%d files", totalFiles)
	}
	title += " in "

	if len(d.changesSinceLastPush) == 1 {
		title += "1 repository"
	} else {
		title += fmt.Sprintf("%d repositories", len(d.changesSinceLastPush))
	}

	description := ""
	for remote, count := range d.changesSinceLastPush {
		if count == 1 {
			description += fmt.Sprintf(" - %s: Upto 1 file\n", remote)
		} else {
			description += fmt.Sprintf(" - %s: Upto %d files\n", remote, count)
		}
	}

	cmd := exec.Command("git", "commit", "-m", title, "-m", description)
	cmd.Dir = d.Directory
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit changes: %s", err)
	}

	d.changesSinceLastPush = make(map[string]int)

	cmd = exec.Command("git", "push")
	cmd.Dir = d.Directory
	return cmd.Run()
}

func (d *Repo) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = d.Directory
	return cmd.Run()
}
