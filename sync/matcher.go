package sync

import (
	"io/fs"
	"path/filepath"

	"github.com/gobwas/glob"

	"github.com/ofabel/fssdk/contract"
)

func ListFiles(root string, include []string, exclude []string, handler contract.FileWalker) error {
	includes := make([]glob.Glob, len(include))

	var err error

	for i, inc := range include {
		if includes[i], err = glob.Compile(inc); err != nil {
			return err
		}
	}

	excludes := make([]glob.Glob, len(exclude))

	for i, exc := range exclude {
		if excludes[i], err = glob.Compile(exc); err != nil {
			return err
		}
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			for _, exc := range excludes {
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

		for _, inc := range includes {
			if inc.Match(path) {
				use = true

				break
			}
		}

		for _, exc := range excludes {
			if exc.Match(path) {
				use = false

				break
			}
		}

		if use {
			full_path := filepath.Join(root, path)

			file := &contract.File{
				Name: filepath.Base(path),
				Path: filepath.Clean(full_path),
				Size: info.Size(),
			}

			if err := handler(file); err != nil {
				return err
			}
		}

		return nil
	})
}
