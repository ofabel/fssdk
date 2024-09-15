package main

import (
	"fmt"
	"log"

	"github.com/ofabel/fssdk/flipper"
)

func main() {
	c, err := flipper.GetConfigFromFile("flipper.json")

	if err != nil {
		log.Fatal(err)
	}

	println(c.Source)

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

	err = f0.UploadFile("main.go", "/ext/test/huge.jpg", func(progress float32) {
		fmt.Printf("%d%%\r", int(progress*100))
	})

	if err != nil {
		log.Fatal(err)
	}

	files, err := f0.GetTree("/ext/test")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		println(file.Path)
	}

	f0.Close()
}
