package persistance

import (
	"fmt"
	"net/url"
	"os/exec"
	"private-ide-config-sync/fs"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
)

var DefaultDatabaseDir string

func init() {
	DefaultDatabaseDir = fmt.Sprintf("%s/.ide-config-sync", fs.HomeDir())
}

func originURLToDir(originURL url.URL) string {
	return fmt.Sprintf("%s%s", originURL.Host, originURL.Path)
}

type IDEConfig struct {
	FsPath       string
	RelativePath string
}

type DatabaseRepo struct {
	directory string
	repo      *git.Repository
}

func NewDatabase(directory string) (*DatabaseRepo, error) {
	repo, err := git.PlainOpen(directory)
	if err != nil {
		return nil, err
	}

	return &DatabaseRepo{
		directory: directory,
		repo:      repo,
	}, nil
}

func (d *DatabaseRepo) Write(origin url.URL, localRepoDir, localFolderPath string) error {
	originDir := originURLToDir(origin)
	originFolderPath := fmt.Sprintf("%s/%s", originDir, localFolderPath)
	dbDir := fmt.Sprintf("%s/%s", d.directory, originDir)
	dbFolderPath := fmt.Sprintf("%s/%s", dbDir, localFolderPath)
	folderPath := fmt.Sprintf("%s%s", localRepoDir, localFolderPath)

	err := copy.Copy(folderPath, dbFolderPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "add", originFolderPath, "--force")
	cmd.Dir = d.directory
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update %s", originDir))
	cmd.Dir = d.directory
	_ = cmd.Run()

	return nil
}

func (d *DatabaseRepo) Read(origin url.URL) ([]IDEConfig, error) {
	originDir := originURLToDir(origin)
	dbDir := fmt.Sprintf("%s/%s", d.directory, originDir)

	exists, _ := fs.Exists(dbDir)
	if !exists {
		return nil, nil
	}

	dirs, err := fs.ListDirs(dbDir)
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
	cmd.Dir = d.directory
	return cmd.Run()
}

func (d *DatabaseRepo) Pull() error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = d.directory
	return cmd.Run()
}
