package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cmd/cli"
	"github.com/ofabel/fssdk/cmd/run"
	"github.com/ofabel/fssdk/cmd/sync"
)

type Args struct {
	Config string     `arg:"-c,--config" help:"Path to the config file." default:"flipper.json"`
	Quiet  bool       `arg:"-q,--quiet" help:"Don't print any output." default:"false"`
	Port   string     `arg:"-p,--port" help:"The port where your Flipper is connected."`
	Cli    *cli.Args  `arg:"subcommand:cli"`
	Run    *run.Args  `arg:"subcommand:run"`
	Sync   *sync.Args `arg:"subcommand:sync"`
}

func (Args) Version() string {
	return fmt.Sprintf("Flipper Zero Script SDK - %s", version)
}

func (Args) Epilogue() string {
	return "For more details visit https://github.com/ofabel/fssdk"
}

func main() {
	var args Args

	parser := arg.MustParse(&args)

	runtime := app.New(args.Config, args.Quiet, args.Port)

	defer runtime.Destroy()

	if cmd := parser.Subcommand(); cmd == nil {
		sync.Main(runtime, nil)
		run.Main(runtime, nil)
	} else if cmd == args.Cli {
		cli.Main(runtime, args.Cli)
	} else if cmd == args.Run {
		run.Main(runtime, args.Run)
	} else if cmd == args.Sync {
		sync.Main(runtime, args.Sync)
	} else {
		panic("unknown subcommand")
	}

	os.Exit(0)
}
