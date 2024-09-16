package app

import (
	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
)

type ProgressHandler func(source *contract.File, target *contract.File, progress float32)

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

	if err != nil {
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

		local_rel_path := file.Source.Rel

		if _, ok := dirs[local_rel_path]; !ok {
			remote_path := base.CleanFlipperPath(target + contract.DirSeparator + local_rel_path)

			dirs[local_rel_path] = remote_path

			if err := rpc.Storage_CreateFolderRecursive(remote_path); err != nil {
				return err
			}
		}

		target_path := base.CleanFlipperPath(target + contract.DirSeparator + local_rel_path)

		err := rpc.Storage_UploadFile(file.Source.Path, target_path, func(progress float32) {
			on_progress(file.Source, file.Target, progress)
		})

		if err != nil {
			return err
		}
	}

	return nil
}
