package app

import (
	"fmt"

	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc"
)

type ProgressHandler func(source string, target string, progress float32)

type SyncStatus uint64

const (
	SyncStatus_Local SyncStatus = iota
	SyncStatus_Both
	SyncStatus_Orphan
)

type SyncFile struct {
	Status SyncStatus
	Source *contract.File
	Target *contract.File
}

type SyncMap map[string]*SyncFile

func (f0 *Flipper) GetSyncMap(source string, target string, includes []string, excludes []string) (SyncMap, error) {
	sync_map := make(SyncMap)

	err := base.ListFiles(source, includes, excludes, func(file *contract.File) error {
		sync_map[file.Rel] = &SyncFile{
			Source: file,
			Status: SyncStatus_Local,
		}

		return nil
	})

	if err != nil {
		return sync_map, err
	}

	files, err := f0.rpc.Storage_GetTree(target)

	if err == rpc.ErrStorageNotExist {
		return sync_map, nil
	} else if err != nil {
		return sync_map, err
	}

	for _, file := range files {
		if sync_file, ok := sync_map[file.Rel]; ok {
			sync_file.Status = SyncStatus_Both
			sync_file.Target = file
		} else {
			sync_map[file.Rel] = &SyncFile{
				Target: file,
				Status: SyncStatus_Orphan,
			}
		}
	}

	return sync_map, err
}

func (f0 *Flipper) SyncFiles(files SyncMap, target string, on_progress ProgressHandler) error {
	rpc, err := f0.GetRpcSession()

	if err != nil {
		return err
	}

	dirs := make(map[string]string)

	for _, file := range files {
		if file.Status == SyncStatus_Orphan {
			continue
		}

		source_dir_path := file.Source.Dir

		if _, ok := dirs[source_dir_path]; !ok {
			target_dir_path := base.Flipper_GetCleanPath(target, source_dir_path)

			dirs[source_dir_path] = target_dir_path

			if false {
				fmt.Printf("mkdir %s\n", target_dir_path)
			} else if err := rpc.Storage_CreateFolderRecursive(target_dir_path); err != nil {
				return err
			}
		}

		target_file_path := base.Flipper_GetCleanPath(target, file.Source.Rel)

		if false {
			fmt.Printf("upload %s\n", target_file_path)

			continue
		}

		err := rpc.Storage_UploadFile(file.Source.Path, target_file_path, func(progress float32) {
			on_progress(file.Source.Path, target_file_path, progress)
		})

		if err != nil {
			return err
		}
	}

	return nil
}
