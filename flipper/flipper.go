package flipper

import (
	"errors"

	"github.com/albenik/go-serial/v2"
)

var ErrNoTerminalFound = errors.New("no terminal found")
var ErrCommandNotSend = errors.New("unable to send command")

type Flipper struct {
	port *serial.Port
	seq  uint32
	rpc  bool
}

var TerminalDelimiter = []byte("\n>: ")
var CRLF = []byte("\r\n")

func Open(port string) (*Flipper, error) {
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

	return &Flipper{
		port: connection,
		seq:  1,
		rpc:  false,
	}, nil
}

func (f0 *Flipper) Close() error {
	return f0.port.Close()
}

func (f0 *Flipper) ReadUntilTerminal() ([]byte, error) {
	data, found, err := f0.ReadUntil(TerminalDelimiter)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrNoTerminalFound
	}

	return data, nil
}

func (f0 *Flipper) ReadUntil(needle []byte) ([]byte, bool, error) {
	character := make([]byte, 1)
	buffer := make([]byte, 1)
	i := 0

	for {
		if i == len(needle) {
			return buffer[1:], true, nil
		}

		n, err := f0.port.Read(character)

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

func (f0 *Flipper) ReadLine() ([]byte, error) {
	character := make([]byte, 1)
	buffer := make([]byte, 1)
	cr := false

	for {
		n, err := f0.port.Read(character)

		if err != nil {
			return nil, err
		}

		if n == 0 {
			return buffer[1:], nil
		}

		if character[0] == '\r' {
			cr = true

			continue
		}

		if cr && character[0] == '\n' {
			return buffer[1:], nil
		}

		cr = false

		buffer = append(buffer, character...)
	}
}

func (f0 *Flipper) SendCommand(command string) error {
	raw_command := []byte(command + "\r")

	n, err := f0.port.Write(raw_command)

	if err != nil {
		return err
	}

	if n != len(raw_command) {
		return ErrCommandNotSend
	}

	return nil
}

func (f0 *Flipper) StartRpcSession() error {
	if f0.rpc {
		return nil
	}

	err := f0.SendCommand("start_rpc_session")

	if err != nil {
		return err
	}

	_, found, err := f0.ReadUntil(CRLF)

	if found && err == nil {
		f0.rpc = true
	}

	return err
}
