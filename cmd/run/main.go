package run

import (
	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cli"
)

const Command = "run"

type Args struct {
	DryRun bool `arg:"-d,--dry-run" help:"Do a dry run, don't execute any commands." default:"false"`
}

const CommandCtrlC = "<CTRL+C>"

type ExecCallback func(command string, output []byte, err error) (bool, error)

func Main(runtime *app.Runtime, args *Args) {
	if args == nil {
		args = &Args{}
	}

	config := runtime.Config()

	if args.DryRun {
		for _, command := range config.Run {
			runtime.Printf("%s%s\n", cli.TerminalDelimiter[1:], command)
		}

		return
	}

	session := runtime.CLI()

	RunCommands(session, config.Run, func(command string, output []byte, err error) (bool, error) {
		if output == nil {
			runtime.Printf("%s%s\n", cli.TerminalDelimiter, command)
		} else {
			runtime.Printf("%s", output)
		}

		return true, err
	})
}

func RunCommands(cli *cli.CLI, commands []string, callback ExecCallback) error {
	var err error

	for _, command := range commands {
		if command == CommandCtrlC {
			err = cli.SendCtrlC()
		} else {
			err = cli.SendCommand(command)
		}

		if next, err := callback(command, nil, err); !next {
			return err
		}

		out, err := cli.ReadUntilTerminal()
		if next, err := callback(command, out, err); !next {
			return err
		}
	}

	return nil
}
