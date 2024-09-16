package app

import (
	"fmt"

	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
)

type SyncStatus uint64

const (
	SyncStatus_Local SyncStatus = iota
	SyncStatus_Both
	SyncStatus_Orphan
)

type SyncMap map[string]SyncStatus

func (f0 *Flipper) GetSyncMap(source string, target string, includes []string, excludes []string) (SyncMap, error) {
	sync_map := make(SyncMap)

	err := base.ListFiles(source, includes, excludes, func(file *contract.File) error {
		sync_map[file.Rel] = SyncStatus_Local

		return nil
	})

	if err != nil {
		return sync_map, err
	}

	files, err := f0.rpc.Storage_GetTree(target)

	if err != nil {
		return sync_map, err
	}

	for _, file := range files {
		if _, ok := sync_map[file.Rel]; ok {
			sync_map[file.Rel] = SyncStatus_Both
		} else {
			sync_map[file.Rel] = SyncStatus_Orphan
		}
	}

	return sync_map, err
}

func (f0 *Flipper) SyncFiles(files []*contract.File, target string) error {
	rpc, err := f0.GetRpcSession()

	if err != nil {
		return err
	}

	dirs := make(map[string]string)

	for _, file := range files {
		if _, ok := dirs[file.Dir]; !ok {
			path := base.CleanFlipperPath(target + contract.DirSeparator + file.Dir)

			dirs[file.Dir] = path

			if err := rpc.Storage_CreateFolderRecursive(path); err != nil {
				return err
			}
		}

		path := base.CleanFlipperPath(target + contract.DirSeparator + file.Path)

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
