package ide

var ideFolderPatterns = []string{
	".idea",
	"\\.vscode",
}

var ignorePatterns = []string{
	"node_modules",
	"venv",
	"\\.venv",
}

type Config struct {
	FsPath       string
	RelativePath string
}
