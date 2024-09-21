package sync

import (
	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc"
)

type UploadProgressState int

const (
	UploadProgressState_DryRun UploadProgressState = 1 << iota
	UploadProgressState_Upload
	UploadProgressState_Skip
)

type ProgressHandler func(state UploadProgressState, source string, target string, progress float32)
type MakeFolderHandler func(dry_run bool, source string, target string)

type SyncStatus string

const (
	SyncStatus_Local  SyncStatus = "< "
	SyncStatus_Both   SyncStatus = "<>"
	SyncStatus_Orphan SyncStatus = " >"
)

type SyncFile struct {
	Status SyncStatus
	Source *contract.File
	Target *contract.File
}

type SyncMap map[string]*SyncFile

func GetLocalSyncMap(path string, includes []string, excludes []string) (SyncMap, error) {
	sync_map := make(SyncMap)

	err := base.ListFiles(path, includes, excludes, func(file *contract.File) error {
		sync_map[file.Rel] = &SyncFile{
			Source: file,
			Status: SyncStatus_Local,
		}

		return nil
	})

	return sync_map, err
}

func GetSyncMap(session *rpc.RPC, source string, target string, includes []string, excludes []string) (SyncMap, error) {
	sync_map, err := GetLocalSyncMap(source, includes, excludes)

	if err != nil {
		return sync_map, err
	}

	files, err := session.Storage_GetTree(target)

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

func SyncFiles(session *rpc.RPC, files SyncMap, target string, dry_run bool, on_progress ProgressHandler, on_make_folder MakeFolderHandler) error {
	dirs := make(map[string]string)

	for _, file := range files {
		if file.Status == SyncStatus_Orphan {
			continue
		}

		source_dir_path := file.Source.Dir

		if _, ok := dirs[source_dir_path]; !ok {
			target_dir_path := base.Flipper_GetCleanPath(target, source_dir_path)

			dirs[source_dir_path] = target_dir_path

			var err error
			var found bool

			if dry_run {
				found, err = session.Storage_FolderExists(target_dir_path)
			}

			if err != nil {
				return err
			} else if dry_run && !found {
				on_make_folder(true, source_dir_path, target_dir_path)
			} else if created, err := session.Storage_CreateFolderRecursive(target_dir_path); err != nil {
				return err
			} else if created {
				on_make_folder(false, source_dir_path, target_dir_path)
			}
		}

		target_file_path := base.Flipper_GetCleanPath(target, file.Source.Rel)

		var same bool

		if dry_run {
			same, _ = session.Storage_CheckFilesAreSame(file.Source.Path, target_file_path)
		}

		if !dry_run {
			// NOP
		} else if same {
			on_progress(UploadProgressState_DryRun|UploadProgressState_Skip, file.Source.Path, target_file_path, 1)

			continue
		} else {
			on_progress(UploadProgressState_DryRun|UploadProgressState_Upload, file.Source.Path, target_file_path, 1)

			continue
		}

		err := session.Storage_UploadFile(file.Source.Path, target_file_path, func(skip bool, progress float32) {
			status := UploadProgressState_Upload

			if skip {
				status = UploadProgressState_Skip
			}

			on_progress(status, file.Source.Path, target_file_path, progress)
		})

		if err != nil {
			return err
		}
	}

	return nil
}
