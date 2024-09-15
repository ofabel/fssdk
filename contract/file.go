package contract

const DirSeparator = "/"
const ThisDirectory = "."
const ExtStorageBasePath = "/ext"

type File struct {
	Name string
	Path string
	Dir  string
	Size int64
}

type FileWalker func(file *File) error
