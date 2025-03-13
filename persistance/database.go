package persistance

import (
	"fmt"
	"ide-config-sync/fsutil"
	"ide-config-sync/ide"
	"ide-config-sync/repository"
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

func (d *DatabaseRepo) writeRemote(remote string, repo, localFolderPath string) error {
	remoteAsPath := remoteURLToDir(remote)
	configFolderPath := fmt.Sprintf("%s/%s", remoteAsPath, localFolderPath)
	dbFolderPath := fmt.Sprintf("%s/%s/%s", d.Directory, remoteAsPath, localFolderPath)
	localFolderAbsPath := fmt.Sprintf("%s/%s", repo, localFolderPath)

	err := copy.Copy(localFolderAbsPath, dbFolderPath)
	if err != nil {
		return fmt.Errorf("failed to copy %s to %s: %s", localFolderAbsPath, dbFolderPath, err)
	}

	cmd := exec.Command("git", "add", configFolderPath, "--force")
	cmd.Dir = d.Directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add %s to git: %s", configFolderPath, err)
	}

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update '%s' with '%s'", remoteAsPath, localFolderPath))
	cmd.Dir = d.Directory
	_ = cmd.Run()

	return nil
}

func (d *DatabaseRepo) Write(repo, relativeConfigPath string) error {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		err := d.writeRemote(remote, repo, relativeConfigPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DatabaseRepo) readRemote(remote string) ([]ide.Config, error) {
	remoteAsPath := remoteURLToDir(remote)
	dbDir := fmt.Sprintf("%s/%s", d.Directory, remoteAsPath)

	exists, _ := fsutil.Exists(dbDir)
	if !exists {
		return nil, nil
	}

	dirs := ide.ReadIDEFolderPaths(dbDir)

	ideConfigs := make([]ide.Config, 0)
	for dir := range dirs {
		ideConfigs = append(ideConfigs, ide.Config{
			FsPath:       fmt.Sprintf("%s/%s", dbDir, dir),
			RelativePath: strings.TrimPrefix(dir, "/"),
		})
	}

	return ideConfigs, nil
}

func (d *DatabaseRepo) Read(repo string) ([]ide.Config, error) {
	remotes, err := repository.ReadRemotes(repo)
	if err != nil {
		return nil, err
	}

	var ideConfigs []ide.Config
	for _, remote := range remotes {
		configs, err := d.readRemote(remote)
		if err != nil {
			return nil, err
		}
		ideConfigs = append(ideConfigs, configs...)
	}
	return ideConfigs, nil
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
