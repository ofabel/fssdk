package contract

const DirSeparator = "/"
const ThisDirectory = "."
const ExtStorageBasePath = "/ext"
const IntStorageBasePath = "/int"

type File struct {
	Name string
	Path string
	Dir  string
	Rel  string
	Size int64
}

type FileWalker func(file *File) error
