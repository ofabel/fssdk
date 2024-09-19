package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cli"
)

type Args struct {
	Command string `arg:"-C,--command" help:"Execute a single command."`
}

func Main(runtime *app.Runtime, args *Args) {
	session := runtime.Terminal()

	defer session.Close()

	if out, err := session.ReadUntilTerminal(); err != nil {
		panic(err)
	} else {
		fmt.Printf("%s%s", out, cli.TerminalDelimiter)
	}

	buffer := make([]byte, 1)
	cmd := ""

	for {
		if n, err := os.Stdin.Read(buffer); err != nil {
			panic(err)
		} else if n == 0 {
			panic("empty input")
		} else if n > 0 {
			cmd = fmt.Sprintf("%s%s", cmd, buffer[0:n])

			if strings.HasSuffix(cmd, "\n") {
				cmd = strings.TrimRight(cmd, "\r\n")

				session.SendCommand(cmd)

				if out, err := session.ReadUntilTerminal(); err != nil {
					panic(err)
				} else {
					fmt.Printf("%s%s", out[len(cmd):], cli.TerminalDelimiter)
				}

				cmd = ""
			}
		}
	}
}
