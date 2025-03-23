package persistance

import (
	"fmt"
	"ide-config-sync/repository"
	"ide-config-sync/util/fsutil"
	"ide-config-sync/util/sliceutil"
	"os/exec"
	"regexp"
	"strings"

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

type DatabaseRepo struct {
	Directory string
	repo      *git.Repository
}

func NewDatabase(directory string) (*DatabaseRepo, error) {
	_, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}

	return &DatabaseRepo{
		Directory: directory,
	}, nil
}

func InitializeFromURL(url, directory string) (*DatabaseRepo, error) {
	cmd := exec.Command("git", "clone", url, directory)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return NewDatabase(directory)
}

func InitializeFromPath(directory string) (*DatabaseRepo, error) {
	_, err := git.PlainInit(directory, false)
	if err != nil {
		return nil, err
	}

	return NewDatabase(directory)
}

func (d *DatabaseRepo) remoteFilesDir(remote string) string {
	remoteAsPath := remoteURLToDir(remote)
	return fmt.Sprintf("%s/%s/files", d.Directory, remoteAsPath)
}

func (d *DatabaseRepo) remoteIncludesFile(remote string) string {
	remoteAsPath := remoteURLToDir(remote)
	return fmt.Sprintf("%s/%s/includes", d.Directory, remoteAsPath)
}

func (d *DatabaseRepo) writeRemoteRepoFile(remote string, repo, localFileRelPath string) error {
	relativePath := fmt.Sprintf("%s/%s", remoteURLToDir(remote), localFileRelPath)

	dbFilePath := fmt.Sprintf("%s/%s", d.remoteFilesDir(remote), relativePath)
	localFileAbsPath := fmt.Sprintf("%s/%s", repo, localFileRelPath)

	err := copy.Copy(localFileAbsPath, dbFilePath)
	if err != nil {
		return fmt.Errorf("failed to copy %s to %s: %s", localFileAbsPath, dbFilePath, err)
	}

	cmd := exec.Command("git", "add", relativePath, "--force")
	cmd.Dir = d.Directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add %s to git: %s", relativePath, err)
	}

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update '%s' with '%s'", remote, localFileRelPath))
	cmd.Dir = d.Directory
	_ = cmd.Run()

	return nil
}

func (d *DatabaseRepo) WriteRepoFile(repo, relativeFilePath string) error {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		err := d.writeRemoteRepoFile(remote, repo, relativeFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DatabaseRepo) WriteRepoIncludes(repo string, includes []string) error {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return err
	}

	for _, remote := range remotes {
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
	}

	return nil
}

func (d *DatabaseRepo) ReadRepoFiles(repo string) (files []string, err error) {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return nil, err
	}

	for _, remote := range remotes {
		filesPath := d.remoteFilesDir(remote)
		exists, _ := fsutil.Exists(filesPath)
		if !exists {
			continue
		}

		files = append(files, filesPath)
	}

	return files, nil
}

func (d *DatabaseRepo) ReadRepoIncludes(repo string) (includes []string, err error) {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return nil, err
	}

	for _, remote := range remotes {
		includesPath := d.remoteIncludesFile(remote)
		exists, _ := fsutil.Exists(includesPath)
		if !exists {
			continue
		}

		lines, err := fsutil.ReadFileLines(includesPath)
		if err != nil {
			return nil, err
		}

		includes = append(includes, lines...)
	}

	return sliceutil.Unique(includes), nil
}

func (d *DatabaseRepo) Push() error {
	cmd := exec.Command("git", "push")
	cmd.Dir = d.Directory
	return cmd.Run()
}

func (d *DatabaseRepo) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = d.Directory
	return cmd.Run()
}
