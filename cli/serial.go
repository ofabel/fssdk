package cli

import (
	"errors"

	"github.com/albenik/go-serial/v2"
)

type CLI struct {
	port *serial.Port
}

var ErrNoTerminalFound = errors.New("no terminal found")
var ErrCommandNotSend = errors.New("unable to send command")

var TerminalDelimiter = []byte("\n>: ")
var CR = byte('\r')
var LF = byte('\n')
var CRLF = []byte("\r\n")
var CTRL_C = []byte("\x03")

func Open(port string) (*CLI, error) {
	connection, err := serial.Open(port,
		serial.WithBaudrate(230400),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(1000),
		serial.WithWriteTimeout(1000),
		serial.WithHUPCL(true))

	if err != nil {
		return nil, err
	}

	return &CLI{
		port: connection,
	}, nil
}

func (cli *CLI) Close() error {
	return cli.port.Close()
}

func (cli *CLI) Write(data []byte) (int, error) {
	return cli.port.Write(data)
}

func (cli *CLI) Read(data []byte) (int, error) {
	return cli.port.Read(data)
}

func (cli *CLI) ReadUntilTerminal() ([]byte, error) {
	data, found, err := cli.ReadUntil(TerminalDelimiter)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrNoTerminalFound
	}

	return data, nil
}

func (cli *CLI) ReadUntil(needle []byte) ([]byte, bool, error) {
	character := make([]byte, 1)
	buffer := make([]byte, 1)
	i := 0

	for {
		if i == len(needle) {
			return buffer[1:], true, nil
		}

		n, err := cli.port.Read(character)

		if err != nil {
			return nil, false, err
		}

		if n == 0 {
			return buffer[1:], false, nil
		}

		if character[0] == needle[i] {
			i++
		} else {
			i = 0
		}

		buffer = append(buffer, character...)
	}
}

func (cli *CLI) ReadLine() ([]byte, error) {
	character := make([]byte, 1)
	buffer := make([]byte, 1)
	cr := false

	for {
		n, err := cli.port.Read(character)

		if err != nil {
			return nil, err
		}

		if n == 0 {
			return buffer[1:], nil
		}

		if character[0] == CR {
			cr = true

			continue
		}

		if cr && character[0] == LF {
			return buffer[1:], nil
		}

		cr = false

		buffer = append(buffer, character...)
	}
}

func (cli *CLI) SendCommand(command string) error {
	raw_command := []byte(command)

	raw_command = append(raw_command, CR)

	n, err := cli.port.Write(raw_command)

	if err != nil {
		return err
	}

	if n != len(raw_command) {
		return ErrCommandNotSend
	}

	return nil
}

func (cli *CLI) SendCtrlC() error {
	_, err := cli.port.Write(CTRL_C)

	return err
}
