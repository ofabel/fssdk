package main

import (
	"fmt"
	"log"

	"github.com/ofabel/fssdk/flipper"
)

func main() {
	port, err := flipper.GetFlipperPort()

	if err != nil {
		fmt.Println(err)

		return
	}

	f0, err := flipper.Open(port)

	if err != nil {
		print(err)

		return
	}

	f0.ReadUntilTerminal()
	err = f0.StartRpcSession()

	if err != nil {
		log.Fatal(err)
	}

	files := make([]*flipper.File, 0, 128)

	err = f0.WalkStorageFiles("/ext/test", func(file *flipper.File) {
		files = append(files, file)
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		println(file.Path)
	}

	f0.Close()
}
