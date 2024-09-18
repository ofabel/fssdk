package run

import "github.com/ofabel/fssdk/app"

const Command = "run"

type Args struct {
	DryRun bool `arg:"-d,--dry-run" help:"Do a dry run, don't execute any commands." default:"false"`
}

func Main(runtime *app.Runtime, args *Args) {
}
