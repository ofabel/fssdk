package sync

const Command = "sync"

type Args struct {
	DryRun bool `arg:"-d,--dry-run" help:"Do a dry run, don't upload, download or delete any files." default:"false"`
	Force  bool `arg:"-f,--force" help:"Upload without checks." default:"false"`
}

func Main(args *Args) {
}
