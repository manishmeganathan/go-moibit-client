package moibit

type FilePath struct {
	elements  []string
	extension string
}

func Root() FilePath {
	return FilePath{}
}

func NewFilePath(elements ...string) (FilePath, error) {
	return FilePath{}, nil
}

func NewFilePathFromString(path string) (FilePath, error) {
	return FilePath{}, nil
}

func (fp FilePath) Path() string {
	return ""
}

func (fp FilePath) IsRoot() bool {
	return fp.IsDirectory() && len(fp.elements) == 0
}

func (fp FilePath) IsDirectory() bool {
	return fp.extension == ""
}

func (fp FilePath) IsFile() bool {
	return fp.extension != ""
}

func (fp *FilePath) Grow(elements ...string) error {
	return nil
}

func (fp FilePath) Parent() FilePath {
	return FilePath{}
}
