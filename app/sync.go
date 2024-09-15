package app

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gobwas/glob"

	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
)

func (f0 *Flipper) ListFiles(root string, include []string, exclude []string, handler contract.FileWalker) error {
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
			dir_path := filepath.Dir(full_path)

			file := &contract.File{
				Name: filepath.Base(path),
				Path: filepath.Clean(full_path),
				Dir:  filepath.Clean(dir_path),
				Size: info.Size(),
			}

			if err := handler(file); err != nil {
				return err
			}
		}

		return nil
	})
}

func (f0 *Flipper) SyncFiles(files []*contract.File, target string) error {
	rpc, err := f0.GetRpcSession()

	if err != nil {
		return err
	}

	dirs := make(map[string]string)

	for _, file := range files {
		if _, ok := dirs[file.Dir]; !ok {
			path := base.CleanPath(target + contract.DirSeparator + file.Dir)

			dirs[file.Dir] = path

			if err := rpc.Storage_CreateFolderRecursive(path); err != nil {
				return err
			}
		}

		path := base.CleanPath(target + contract.DirSeparator + file.Path)

		err := rpc.Storage_UploadFile(file.Path, path, func(progress float32) error {
			fmt.Printf("%s [%d%%]\r", path, int(progress*100))

			return nil
		})

		if err != nil {
			return err
		}

		println(path + "       ")
	}

	return nil
}
