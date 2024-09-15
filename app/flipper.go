package app

import (
	"github.com/ofabel/fssdk/cli"
	"github.com/ofabel/fssdk/rpc"
)

type Flipper struct {
	port string
	cli  *cli.CLI
	rpc  *rpc.RPC
}

func New(port string) *Flipper {
	return &Flipper{
		port: port,
		cli:  nil,
		rpc:  nil,
	}
}

func (f0 *Flipper) StartCliSession() (*cli.CLI, error) {
	if f0.cli != nil {
		return f0.cli, nil
	}

	cli, err := cli.Open(f0.port)

	if err != nil {
		return nil, err
	}

	f0.cli = cli

	return cli, nil
}

func (f0 *Flipper) StartRpcSession() (*rpc.RPC, error) {
	if f0.rpc != nil {
		return f0.rpc, nil
	}

	if _, err := f0.StartCliSession(); err != nil {
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
