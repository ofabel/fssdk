package app

import (
	"errors"

	"github.com/ofabel/fssdk/cli"
	"github.com/ofabel/fssdk/rpc"
)

var ErrRpcSessionActive = errors.New("RPC session is active")

type Flipper struct {
	port   string
	config *Config
	cli    *cli.CLI
	rpc    *rpc.RPC
}

func New(port string, config *Config) *Flipper {
	return &Flipper{
		port:   port,
		config: config,
		cli:    nil,
		rpc:    nil,
	}
}

func (f0 *Flipper) Close() error {
	if f0.cli != nil {
		return f0.cli.Close()
	} else {
		return nil
	}
}

func (f0 *Flipper) GetCliSession() (*cli.CLI, error) {
	if f0.rpc != nil {
		return nil, ErrRpcSessionActive
	}

	if f0.cli != nil {
		return f0.cli, nil
	}

	var err error

	f0.cli, err = cli.Open(f0.port)

	if err != nil {
		return nil, err
	}

	if _, err := f0.cli.ReadUntilTerminal(); err != nil {
		return nil, err
	}

	return f0.cli, nil
}

func (f0 *Flipper) GetRpcSession() (*rpc.RPC, error) {
	if f0.rpc != nil {
		return f0.rpc, nil
	}

	if _, err := f0.GetCliSession(); err != nil {
		return nil, err
	}

	err := f0.cli.SendCommand("start_rpc_session")

	if err != nil {
		return nil, err
	}

	_, found, err := f0.cli.ReadUntil(cli.CRLF)

	if found && err == nil {
		f0.rpc = rpc.New(f0.cli)
	}

	return f0.rpc, err
}

func (f0 *Flipper) StopRpcSession() error {
	if f0.rpc == nil {
		return nil
	}

	if err := f0.rpc.StopSession(); err != nil {
		return err
	}

	f0.rpc = nil

	cli, err := f0.GetCliSession()

	if err != nil {
		return err
	}

	_, err = cli.ReadUntilTerminal()

	return err
}
