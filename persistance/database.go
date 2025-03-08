package persistance

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
)

var DefaultDatabaseDir string

func init() {
	DefaultDatabaseDir = filepath.ToSlash(fmt.Sprintf("%s/.ide-config-sync", HomeDir()))
}

func originURLToDir(originURL url.URL) string {
	return fmt.Sprintf("%s%s", originURL.Host, originURL.Path)
}

type IDEConfig struct {
	FsPath       string
	RelativePath string
}

type DatabaseRepo struct {
	Directory string
	repo      *git.Repository
}

func NewDatabase(directory string) (*DatabaseRepo, error) {
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}

	return &DatabaseRepo{
		Directory: directory,
		repo:      repo,
	}, nil
}

func NewDatabaseFromURL(url, directory string) (*DatabaseRepo, error) {
	cmd := exec.Command("git", "clone", url, directory)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}

	return &DatabaseRepo{
		Directory: directory,
		repo:      repo,
	}, nil
}

func (d *DatabaseRepo) Write(origin url.URL, localRepoDir, localFolderPath string) error {
	originDir := originURLToDir(origin)
	originFolderPath := fmt.Sprintf("%s/%s", originDir, localFolderPath)
	dbFolderPath := fmt.Sprintf("%s/%s/%s", d.Directory, originDir, localFolderPath)
	localFolderAbsPath := fmt.Sprintf("%s%s", localRepoDir, localFolderPath)

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

func (d *DatabaseRepo) Read(origin url.URL) ([]IDEConfig, error) {
	originDir := originURLToDir(origin)
	dbDir := fmt.Sprintf("%s/%s", d.Directory, originDir)

	exists, _ := Exists(dbDir)
	if !exists {
		return nil, nil
	}

	dirs, err := ListDirs(dbDir)
	if err != nil {
		return nil, err
	}

	ideConfigs := make([]IDEConfig, 0)
	for _, dir := range dirs {
		relativePath := strings.Split(dir, originDir)[1]
		ideConfigs = append(ideConfigs, IDEConfig{
			FsPath:       dir,
			RelativePath: relativePath,
		})
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
