package database

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"ide-config-sync/repository"
	"ide-config-sync/set"
	"ide-config-sync/util/fsutil"

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
	Directory string
	repo      *git.Repository
}

func NewRepoDatabase(directory string) (Database, error) {
	_, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Directory: directory,
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

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update '%s' with '%s'", remote, file.PathFromRepoRoot))
	cmd.Dir = d.Directory
	_ = cmd.Run()

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

func (d *Repo) WriteRepoOptions(repo string, options repository.Options) error {
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

func (d *Repo) ReadRepoFiles(repo string) ([]repository.File, error) {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return nil, err
	}

	var files []repository.File
	for _, remote := range remotes {
		filesPath := d.remoteFilesDir(remote)
		exists, _ := fsutil.Exists(filesPath)
		if !exists {
			continue
		}

		files = append(files, repository.File{
			Type:             repository.FileTypeDir,
			PathFromRepoRoot: "/",
			AbsolutePath:     filesPath,
		})
	}

	return files, nil
}

func (d *Repo) ReadRepoOptions(repo string) (repository.Options, error) {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return repository.Options{}, err
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
			return repository.Options{}, err
		}

		includes = append(includes, lines...)
	}

	return repository.NewOptions(includes), nil
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

func (d *Repo) ReadGlobalOptions() (repository.Options, error) {
	if !d.globalOptionsExists() {
		options := repository.DefaultGlobalOptions()
		err := d.WriteGlobalOptions(options)
		if err != nil {
			return repository.Options{}, err
		}
		return options, nil
	}
	includeLines, inclErr := d.readGlobalIncludes()
	excludeLines, exclErr := d.readGlobalExcludes()

	if inclErr != nil || exclErr != nil {
		return repository.Options{}, fmt.Errorf("failed to read global includes / excludes")
	}

	options := repository.DefaultGlobalOptions()
	options.IncludePatterns = set.FromSlice(includeLines)
	options.ExcludePatterns = set.FromSlice(excludeLines)

	return options, nil
}

func (d *Repo) WriteGlobalOptions(config repository.Options) error {
	err := fsutil.WriteFileLines(d.globalIncludesFile(), config.IncludePatterns.Slice())
	if err != nil {
		return err
	}

	err = fsutil.WriteFileLines(d.globalExcludesFile(), config.ExcludePatterns.Slice())
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

	cmd = exec.Command("git", "commit", "-m", "Update global includes / excludes")
	cmd.Dir = d.Directory
	_ = cmd.Run()

	return nil
}

func (d *Repo) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = d.Directory
	return cmd.Run()
}

func (d *Repo) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = d.Directory
	return cmd.Run()
}
