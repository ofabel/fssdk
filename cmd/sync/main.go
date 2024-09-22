package sync

import (
	"fmt"

	"github.com/ofabel/fssdk/app"
)

const Command = "sync"

type Args struct {
	DryRun bool   `arg:"-d,--dry-run" help:"Do a dry run, don't upload, download or delete any files." default:"false"`
	Force  bool   `arg:"-f,--force" help:"Upload without checks." default:"false"`
	List   bool   `arg:"-l,--list" help:"List matching files." default:"false"`
	Local  bool   `arg:"-o,--local" help:"List matching files from local source only." default:"false"`
	Source string `arg:"-s,--source" help:"Sync all from source to target. If source is a folder, target is also treated as a folder."`
	Target string `arg:"-t,--target" help:"Sync all from source to target."`
}

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

	if args.Local {
		config := runtime.Config()

		sync_map, err := GetLocalSyncMap(source, config.Include, config.Exclude)

		if err != nil {
			panic(err)
		}

		ListFilesFromSyncMap(runtime, sync_map)

		return
	}

	if args.List {
		config := runtime.Config()
		session := runtime.RPC()

		sync_map, err := GetSyncMap(session, source, target, config.Include, config.Exclude)

		if err != nil {
			panic(err)
		}

		ListFilesFromSyncMap(runtime, sync_map)

		return
	}

	session := runtime.RPC()
	config := runtime.Config()

	files, err := GetSyncMap(session, source, target, config.Include, config.Exclude)

	if err != nil {
		panic(err)
	}

	width := getMaxWidth(files)

	on_progress := func(direction TransferDirection, operation TransferOperation, dry_run bool, source string, target string, progress float32) {
		source = padRight(source, width)

		// skip
		if operation != TransferOperation_Handle {
			status := formatStatus("%s", operation)

			runtime.Printf("%s%s %s %s\n", status, source, direction, target)

			return
		}

		// dry run upload
		if direction == TransferDirection_Upload && dry_run {
			status := formatStatus("UPLD")

			runtime.Printf("%s%s %s %s\n", status, source, direction, target)

			return
		}

		// dry run download
		if direction == TransferDirection_Download && dry_run {
			status := formatStatus("DNLD")

			runtime.Printf("%s%s %s %s\n", status, source, direction, target)

			return
		}

		// upload
		if progress < 1 {
			status := formatStatus("%d%%", int(progress*100))

			runtime.Printf("%s%s %s %s\r", status, source, direction, target)

			return
		} else {
			status := formatStatus("100%%")

			runtime.Printf("%s%s %s %s\n", status, source, direction, target)

			return
		}
	}

	on_make_folder := func(dry_run bool, direction TransferDirection, source string, target string) {
		source = padRight(source, width)

		status := formatStatus("MKFD")

		runtime.Printf("%s%s %s %s\n", status, source, direction, target)
	}

	if err := SyncFiles(session, files, source, target, config.Orphans, args.DryRun, on_progress, on_make_folder); err != nil {
		panic(err)
	}
}

func formatStatus(format string, args ...any) string {
	status := fmt.Sprintf(format, args...)

	for len(status) < 4 {
		status = " " + status
	}

	return fmt.Sprintf("[%s]  ", status)
}

func padRight(str string, width int) string {
	for len(str) < width {
		str += " "
	}

	return str
}
