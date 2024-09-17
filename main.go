package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cli"
	"github.com/ofabel/fssdk/cmd/sync"
)

var args struct {
	Config string     `arg:"-c,--config" help:"Path to the config file." default:"flipper.json"`
	Silent bool       `arg:"-s,--silent" help:"Don't print any output" default:"false"`
	Sync   *sync.Args `arg:"subcommand:sync"`
}

func main() {
	parser := arg.MustParse(&args)

	if cmd := parser.Subcommand(); cmd == nil {
		panic(args.Config)
	} else {
		panic(args.Sync.DryRun)
	}

	config, err := app.GetConfigFromFile(args.Config)

	if err != nil {
		log.Fatal(err)
	}

	port, err := cli.GetFlipperPort()

	if err != nil {
		log.Fatal(err)

		return
	}

	app := app.New(port, config)

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
}
