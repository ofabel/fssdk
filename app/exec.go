package app

const CommandCtrlC = "<CTRL+C>"

type ExecCallback func(command string, output []byte, err error) (bool, error)

func (f0 *Flipper) RunCommands(commands []string, callback ExecCallback) error {
	cli, err := f0.GetCliSession()

	if err != nil {
		return err
	}

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
