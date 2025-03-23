package repository

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
