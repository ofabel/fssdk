package main

import (
	"fmt"
	"log"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cli"
)

func main() {
	config, err := app.GetConfigFromFile("flipper.json")

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
