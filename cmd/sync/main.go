package sync

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc"
	"golang.org/x/exp/maps"
)

const Command = "sync"

type Args struct {
	DryRun bool   `arg:"-d,--dry-run" help:"Do a dry run, don't upload, download or delete any files." default:"false"`
	Force  bool   `arg:"-f,--force" help:"Upload without checks." default:"false"`
	List   bool   `arg:"-l,--list" help:"List matching local files." default:"false"`
	Source string `arg:"-s,--source" help:"Sync all from source to target. If source is a folder, target is also treated as a folder."`
	Target string `arg:"-t,--target" help:"Sync all from source to target."`
}

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

func Main(runtime *app.Runtime, args *Args) {
	source := args.Source
	target := args.Target

	if len(source) == 0 {
		config := runtime.Config()

		source = config.Source
	}

	if len(target) == 0 {
		config := runtime.Config()

		target = config.Target
	}

	source = runtime.GetAbsolutePath(source)

	if args.List {
		config := runtime.Config()

		sync_map, err := GetLocalSyncMap(source, config.Include, config.Exclude)
		files := maps.Keys(sync_map)

		sort.Slice(files, func(i, j int) bool { // TODO: use https://pkg.go.dev/slices#SortStableFunc
			is := strings.Split(files[i], "/")
			js := strings.Split(files[j], "/")
			size := len(is)

			if len(js) < size {
				size = len(js)
			}

			if size == 1 && len(is) != len(js) {
				return false
			}

			for p := range size {
				if is[p] == js[p] {
					continue
				} else {
					return is[p] < js[p]
				}
			}

			return false
		})

		if err != nil {
			panic(err)
		}

		for _, file := range files {
			runtime.Printf("%s\n", file)
		}

		return
	}

	panic("done")

	session := runtime.RPC()

	config := runtime.Config()

	files, err := GetSyncMap(session, source, target, config.Include, config.Exclude)

	if err != nil {
		panic(err)
	}

	if err := SyncFiles(session, files, target, func(source string, target string, progress float32) {
		if progress < 1 {
			fmt.Printf("%s [%d%%]\r", target, int(progress*100))
		} else {
			fmt.Printf("%s [100%%]\n", target)
		}
	}); err != nil {
		panic(err)
	}
}

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

func SyncFiles(session *rpc.RPC, files SyncMap, target string, on_progress ProgressHandler) error {
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
			} else if err := session.Storage_CreateFolderRecursive(target_dir_path); err != nil {
				return err
			}
		}

		target_file_path := base.Flipper_GetCleanPath(target, file.Source.Rel)

		if false {
			fmt.Printf("upload %s\n", target_file_path)

			continue
		}

		err := session.Storage_UploadFile(file.Source.Path, target_file_path, func(progress float32) {
			on_progress(file.Source.Path, target_file_path, progress)
		})

		if err != nil {
			return err
		}
	}

	return nil
}
