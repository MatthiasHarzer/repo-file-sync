package persistance

import (
	"fmt"
	"ide-config-sync/fs"
	"ide-config-sync/repository"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
)

var DefaultDatabaseDir string

func init() {
	DefaultDatabaseDir = filepath.ToSlash(fmt.Sprintf("%s/.ide-config-sync", fs.HomeDir()))
}

func originURLToDir(originURL url.URL) string {
	return fmt.Sprintf("%s%s", originURL.Host, originURL.Path)
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

func (d *DatabaseRepo) writeOrigin(origin url.URL, localRepoDir, localFolderPath string) error {
	originDir := originURLToDir(origin)
	originFolderPath := fmt.Sprintf("%s/%s", originDir, localFolderPath)
	dbFolderPath := fmt.Sprintf("%s/%s/%s", d.Directory, originDir, localFolderPath)
	localFolderAbsPath := fmt.Sprintf("%s/%s", localRepoDir, localFolderPath)

	err := copy.Copy(localFolderAbsPath, dbFolderPath)
	if err != nil {
		return fmt.Errorf("failed to copy %s to %s: %s", localFolderAbsPath, dbFolderPath, err)
	}

	cmd := exec.Command("git", "add", originFolderPath, "--force")
	cmd.Dir = d.Directory
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add %s to git: %s", originFolderPath, err)
	}

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update '%s' with '%s'", originDir, localFolderPath))
	cmd.Dir = d.Directory
	_ = cmd.Run()

	return nil
}

func (d *DatabaseRepo) Write(origins []url.URL, localRepoDir, localFolderPath string) error {
	for _, origin := range origins {
		err := d.writeOrigin(origin, localRepoDir, localFolderPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DatabaseRepo) readOrigin(origin url.URL) ([]repository.IDEConfig, error) {
	originDir := originURLToDir(origin)
	dbDir := fmt.Sprintf("%s/%s", d.Directory, originDir)

	exists, _ := fs.Exists(dbDir)
	if !exists {
		return nil, nil
	}

	dirs := repository.GetIDEFolderPaths(dbDir)

	ideConfigs := make([]repository.IDEConfig, 0)
	for dir := range dirs {
		ideConfigs = append(ideConfigs, repository.IDEConfig{
			FsPath:       fmt.Sprintf("%s/%s", dbDir, dir),
			RelativePath: strings.TrimPrefix(dir, "/"),
		})
	}

	return ideConfigs, nil
}

func (d *DatabaseRepo) Read(origins []url.URL) ([]repository.IDEConfig, error) {
	var ideConfigs []repository.IDEConfig
	for _, origin := range origins {
		configs, err := d.readOrigin(origin)
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
