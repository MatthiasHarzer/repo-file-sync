package ide

import (
	"fmt"

	"github.com/otiai10/copy"
)

func WriteIDEFolder(repo string, config Config) error {
	repoFolderPath := fmt.Sprintf("%s/%s", repo, config.RelativePath)
	err := copy.Copy(config.FsPath, repoFolderPath)
	if err != nil {
		return err
	}
	return nil
}
