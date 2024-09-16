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

	for file, state := range files {
		fmt.Printf("%s : %v\n", file, state.Status)
	}

	if err := app.StopRpcSession(); err != nil {
		log.Fatal(err)
	}
}
