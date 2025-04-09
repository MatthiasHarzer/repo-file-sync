package repository

import "path/filepath"

type File struct {
	AbsolutePath     string
	PathFromRepoRoot string
}

func splitPath(path string) []string {
	dir, file := filepath.Split(path)
	if dir == "" {
		return []string{file}
	}
	return append(splitPath(filepath.Clean(dir)), file)
}

func fixPath(path string) string {
	path = filepath.ToSlash(path)
	path = filepath.Clean(path)

}

func CopyFile(src, dst string) error {
	srcAdjusted := ""
	destAdjusted := ""

}
