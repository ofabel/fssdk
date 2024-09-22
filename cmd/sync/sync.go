package sync

import (
	"errors"
	"os"
	"path/filepath"
	"slices"

	"github.com/gobwas/glob"
	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc"
	"golang.org/x/exp/maps"
)

type TransferDirection string

const (
	TransferDirection_Upload   TransferDirection = ">"
	TransferDirection_Download TransferDirection = "<"
	TransferDirection_None     TransferDirection = ":"
)

type TransferOperation string

const (
	TransferOperation_Handle TransferOperation = "HNDL"
	TransferOperation_Delete TransferOperation = "DELT"
	TransferOperation_Ignore TransferOperation = "IGNR"
	TransferOperation_Skip   TransferOperation = "SKIP"
)

type ProgressHandler func(direction TransferDirection, operation TransferOperation, dry_run bool, source string, target string, progress float32)
type MakeFolderHandler func(dry_run bool, direction TransferDirection, source string, target string)

const NotExistingFile = "-"

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

func GetLocalSyncMap(path string, includes []glob.Glob, excludes []glob.Glob) (SyncMap, error) {
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

func GetSyncMap(session *rpc.RPC, source string, target string, includes []glob.Glob, excludes []glob.Glob) (SyncMap, error) {
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
		} else if base.CanUse(file.Rel, includes, excludes) {
			rel_path := filepath.FromSlash(file.Rel)
			dir_path := filepath.Dir(rel_path)
			abs_path := filepath.Join(source, rel_path)

			sync_map[file.Rel] = &SyncFile{
				Target: file,
				Source: &contract.File{
					Name: file.Name,
					Path: base.CleanPath(abs_path),
					Dir:  base.CleanPath(dir_path),
					Rel:  base.CleanPath(rel_path),
					Size: file.Size,
				},
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
			if err := handleOrphan(session, file, source, target, orphans, dry_run, dirs, on_progress, on_make_folder); err != nil {
				return err
			}

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
				on_make_folder(true, TransferDirection_Upload, source_dir_path, target_dir_path)
			} else if created, err := session.Storage_CreateFolderRecursive(target_dir_path); err != nil {
				return err
			} else if created {
				on_make_folder(false, TransferDirection_Upload, source_dir_path, target_dir_path)
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
			on_progress(TransferDirection_Upload, TransferOperation_Skip, true, file.Source.Path, target_file_path, 1)

			continue
		} else {
			on_progress(TransferDirection_Upload, TransferOperation_Handle, true, file.Source.Path, target_file_path, 1)

			continue
		}

		err := session.Storage_UploadFile(file.Source.Path, target_file_path, func(skip bool, progress float32) {
			operation := TransferOperation_Handle

			if skip {
				operation = TransferOperation_Skip
			}

			on_progress(TransferDirection_Upload, operation, false, file.Source.Path, target_file_path, progress)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func handleOrphan(session *rpc.RPC, file *SyncFile, source string, target string, orphans contract.Orphans, dry_run bool, dirs map[string]string, on_progress ProgressHandler, on_make_folder MakeFolderHandler) error {
	if orphans == contract.Orphans_Ignore {
		on_progress(TransferDirection_None, TransferOperation_Ignore, dry_run, NotExistingFile, file.Target.Path, 1)

		return nil
	}

	if _, ok := dirs[file.Source.Dir]; orphans == contract.Orphans_Download && !ok {
		target_dir_path := base.Flipper_GetCleanPath(target, file.Target.Dir)

		dirs[file.Source.Dir] = target_dir_path

		stat, err := os.Stat(source)

		if err != nil {
			return err
		}

		source_dir_perm := stat.Mode().Perm()
		source_dir_path := filepath.Dir(file.Source.Path)

		if _, err := os.Stat(source_dir_path); os.IsNotExist(err) {
			on_make_folder(dry_run, TransferDirection_Download, source_dir_path, target_dir_path)
		}

		if !dry_run {
			if err := os.MkdirAll(source_dir_path, source_dir_perm); err != nil {
				return err
			}
		}

	}

	if !dry_run {
		// NOP
	} else if orphans == contract.Orphans_Delete {
		on_progress(TransferDirection_Upload, TransferOperation_Delete, true, NotExistingFile, file.Target.Path, 1)

		return nil
	} else if orphans == contract.Orphans_Download {
		on_progress(TransferDirection_Download, TransferOperation_Handle, true, file.Source.Path, file.Target.Path, 1)

		return nil
	}

	if orphans == contract.Orphans_Download {
		local_file_path := filepath.FromSlash(file.Source.Path)

		return session.Storage_DownloadFile(file.Target.Path, local_file_path, func(skip bool, progress float32) {
			operation := TransferOperation_Handle

			if skip {
				operation = TransferOperation_Skip
			}

			on_progress(TransferDirection_Download, operation, false, file.Source.Path, file.Target.Path, progress)
		})
	}

	if orphans == contract.Orphans_Delete {
		on_progress(TransferDirection_Upload, TransferOperation_Delete, false, NotExistingFile, file.Target.Path, 1)

		return session.Storage_Delete(file.Target.Path, false)
	}

	return nil
}
