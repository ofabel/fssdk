package base

import (
	"io/fs"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/ofabel/fssdk/contract"
)

func ListFiles(root string, includes []glob.Glob, excludes []glob.Glob, handler contract.FileWalker) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		path_to_glob := filepath.ToSlash(path)

		if d.IsDir() {
			for _, exc := range excludes {
				if exc.Match(path_to_glob) {
					return filepath.SkipDir
				}
			}

			return nil
		}

		if !d.Type().IsRegular() {
			return nil
		}

		info, err := d.Info()

		if err != nil {
			return err
		}

		if CanUse(path_to_glob, includes, excludes) {
			rel_path, err := filepath.Rel(root, path)

			if err != nil {
				return err
			}

			dir_path := filepath.Dir(rel_path)

			file := &contract.File{
				Name: filepath.Base(path),
				Path: CleanPathSlash(path),
				Dir:  CleanPathSlash(dir_path),
				Rel:  CleanPathSlash(rel_path),
				Size: info.Size(),
			}

			if err := handler(file); err != nil {
				return err
			}
		}

		return nil
	})
}

func CleanPath(path string) string {
	clean_path := filepath.FromSlash(path)

	return filepath.Clean(clean_path)
}

func CleanPathSlash(path string) string {
	clean_path := filepath.Clean(path)

	return filepath.ToSlash(clean_path)
}

func ToGlobArray(filters []string) ([]glob.Glob, error) {
	globs := make([]glob.Glob, len(filters))

	var err error

	for i, filter := range filters {
		if globs[i], err = glob.Compile(filter); err != nil {
			return globs, err
		}
	}

	return globs, nil
}

func CanUse(subject string, includes []glob.Glob, excludes []glob.Glob) bool {
	use := false

	for _, inc := range includes {
		if inc.Match(subject) {
			use = true

			break
		}
	}

	for _, exc := range excludes {
		if exc.Match(subject) {
			use = false

			break
		}
	}

	return use
}
