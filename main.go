package main

import (
	"fmt"
	"log"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cli"
)

func main() {
	c, err := app.GetConfigFromFile("flipper.json")

	if err != nil {
		log.Fatal(err)
	}

	println(c.Source)

	port, err := cli.GetFlipperPort()

	if err != nil {
		log.Fatal(err)

		return
	}

	app := app.New(port)

	cli, err := app.StartCliSession()

	if err != nil {
		log.Fatal(err)

		return
	}

	defer cli.Close()

	cli.ReadUntilTerminal()

	rpc, err := app.StartRpcSession()

	if err != nil {
		log.Fatal(err)

		return
	}

	err = rpc.Storage_UploadFile("main.go", "/ext/test/huge.jpg", func(progress float32) {
		fmt.Printf("%d%%\r", int(progress*100))
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
}
