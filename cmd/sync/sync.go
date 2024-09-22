package sync

import (
	"errors"
	"slices"

	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc"
	"golang.org/x/exp/maps"
)

type TransferDirection string

const (
	TransferDirection_Upload   TransferDirection = ">"
	TransferDirection_Download TransferDirection = "<"
	TransferDirection_None     TransferDirection = "-"
)

type ProgressHandler func(direction TransferDirection, skip bool, dry_run bool, source string, target string, progress float32)
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

var ErrFileUnknown = errors.New("unknown file")

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

func SyncFiles(session *rpc.RPC, files SyncMap, source string, target string, orphans contract.Orphans, dry_run bool, on_progress ProgressHandler, on_make_folder MakeFolderHandler) error {
	dirs := make(map[string]string)
	keys := maps.Keys(files)

	slices.SortStableFunc(keys, compareFiles)
	slices.Reverse(keys)

	for _, key := range keys {
		file, found := files[key]

		if !found {
			return ErrFileUnknown
		}

		if file.Status == SyncStatus_Orphan {
			handleOrphan(session, file, source, target, orphans, dry_run, on_progress, on_make_folder)

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
			on_progress(TransferDirection_Upload, true, true, file.Source.Path, target_file_path, 1)

			continue
		} else {
			on_progress(TransferDirection_Upload, false, true, file.Source.Path, target_file_path, 1)

			continue
		}

		err := session.Storage_UploadFile(file.Source.Path, target_file_path, func(skip bool, progress float32) {
			on_progress(TransferDirection_Upload, skip, false, file.Source.Path, target_file_path, progress)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func handleOrphan(session *rpc.RPC, file *SyncFile, source string, target string, orphans contract.Orphans, dry_run bool, on_progress ProgressHandler, on_make_folder MakeFolderHandler) error {
	if !dry_run {
		// NOP
	} else if orphans == contract.Orphans_Delete {
		on_progress(TransferDirection_Upload, false, true, "", file.Target.Path, 1)
	} else if orphans == contract.Orphans_Download {
		on_progress(TransferDirection_Download, false, true, "", file.Target.Path, 1)
	} else if orphans == contract.Orphans_Ignore {
		on_progress(TransferDirection_None, true, true, "", file.Target.Path, 1)
	}

	return nil
}
