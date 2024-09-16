package base

import (
	"io/fs"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/ofabel/fssdk/contract"
)

func ListFiles(root string, includes []string, excludes []string, handler contract.FileWalker) error {
	include_globs := make([]glob.Glob, len(includes))

	var err error

	for i, inc := range includes {
		if include_globs[i], err = glob.Compile(inc); err != nil {
			return err
		}
	}

	exclude_globs := make([]glob.Glob, len(excludes))

	for i, exc := range excludes {
		if exclude_globs[i], err = glob.Compile(exc); err != nil {
			return err
		}
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			for _, exc := range exclude_globs {
				if exc.Match(path) {
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

		use := false

		for _, inc := range include_globs {
			if inc.Match(path) {
				use = true

				break
			}
		}

		for _, exc := range exclude_globs {
			if exc.Match(path) {
				use = false

				break
			}
		}

		if use {
			rel_path, err := filepath.Rel(root, path)

			if err != nil {
				return err
			}

			dir_path := filepath.Dir(rel_path)

			file := &contract.File{
				Name: filepath.Base(path),
				Path: filepath.Clean(path),
				Dir:  filepath.Clean(dir_path),
				Rel:  filepath.Clean(rel_path),
				Size: info.Size(),
			}

			if err := handler(file); err != nil {
				return err
			}
		}

		return nil
	})
}
