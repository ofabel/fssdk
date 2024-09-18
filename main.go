package main

import (
	"os"

	"github.com/alexflint/go-arg"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cmd/cli"
	"github.com/ofabel/fssdk/cmd/run"
	"github.com/ofabel/fssdk/cmd/sync"
)

var args struct {
	Config string     `arg:"-c,--config" help:"Path to the config file." default:"flipper.json"`
	Quiet  bool       `arg:"-q,--quiet" help:"Don't print any output." default:"false"`
	Port   string     `arg:"-p,--port" help:"The port where your Flipper is connected."`
	Cli    *cli.Args  `arg:"subcommand:cli"`
	Run    *run.Args  `arg:"subcommand:run"`
	Sync   *sync.Args `arg:"subcommand:sync"`
}

func main() {
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
	/*
	   app := app.New("", config)

	   defer app.Close()

	   	if _, err := app.GetRpcSession(); err != nil {
	   		log.Fatal(err)

	   		return
	   	}

	   files, err := app.GetSyncMap(config.Source, config.Target, config.Include, config.Exclude)

	   	if err != nil {
	   		log.Fatal(err)
	   	}

	   	err = app.SyncFiles(files, config.Target, func(source string, target string, progress float32) {
	   		if progress < 1.0 {
	   			fmt.Printf("%s [%d%%]\r", target, int(100*progress))
	   		} else {
	   			fmt.Printf("%s         \n", target)
	   		}
	   	})

	   	if err != nil {
	   		log.Fatal(err)
	   	}

	   	if err := app.StopRpcSession(); err != nil {
	   		log.Fatal(err)
	   	}
	*/
}
