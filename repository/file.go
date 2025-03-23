package repository

import "ide-config-sync/util/fsutil"

type FileType string

const (
	FileTypeFile FileType = "file"
	FileTypeDir  FileType = "dir"
)

type File struct {
	Type             FileType
	AbsolutePath     string
	PathFromRepoRoot string
}

func (f *File) IsFile() bool {
	return f.Type == FileTypeFile
}

func (f *File) ListFiles() ([]string, error) {
	if f.IsFile() {
		return []string{f.AbsolutePath}, nil
	}

	return fsutil.ListFiles(f.AbsolutePath)
}
