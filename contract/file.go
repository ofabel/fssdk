package contract

type File struct {
	Name string
	Path string
	Size int64
}

type FileWalker func(file *File) error
