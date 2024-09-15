package main

import (
	"fmt"
	"log"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cli"
	"github.com/ofabel/fssdk/contract"
)

func main() {
	config, err := app.GetConfigFromFile("flipper.json")

	if err != nil {
		log.Fatal(err)
	}

	println(config.Source)

	port, err := cli.GetFlipperPort()

	if err != nil {
		log.Fatal(err)

		return
	}

	app := app.New(port)

	defer app.Close()

	rpc, err := app.GetRpcSession()

	if err != nil {
		log.Fatal(err)

		return
	}

	err = rpc.Storage_UploadFile("main.go", "/ext/test/huge.jpg", func(progress float32) error {
		fmt.Printf("%d%%\r", int(progress*100))

		return nil
	})

	if err != nil {
		log.Fatal(err)

		return
	}

	files, err := rpc.Storage_GetTree("/ext/test")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		println(file.Path)
	}

	files = make([]*contract.File, 0, 32)

	err = app.ListFiles(config.Source, config.Include, config.Exclude, func(file *contract.File) error {
		files = append(files, file)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	err = app.SyncFiles(files, config.Target)

	if err != nil {
		log.Fatal(err)
	}

	err = rpc.Storage_CreateFolderRecursive("/ext/test/a/b/c/d/e/f")

	if err != nil {
		log.Fatal(err)
	}

	if err := app.StopRpcSession(); err != nil {
		log.Fatal(err)
	}

	err = app.RunCommands(config.Run, func(command string, output []byte, err error) (bool, error) {
		if err != nil && err != cli.ErrNoTerminalFound {
			return false, err
		}

		fmt.Printf("%s\n", output)

		return true, nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
