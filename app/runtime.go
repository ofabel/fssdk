package app

import (
	"path/filepath"

	"github.com/ofabel/fssdk/cli"
	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc"
)

type Runtime struct {
	root_path   string
	config_file string
	quiet       bool
	port        string
	config      *contract.Config
	cli         *cli.CLI
	rpc         *rpc.RPC
}

func New(config string, quiet bool, port string) *Runtime {
	return &Runtime{
		config_file: config,
		quiet:       quiet,
		port:        port,
	}
}

func (r *Runtime) Destroy() {
	if r.rpc != nil {
		r.rpc.Close()
	}

	if r.cli != nil {
		r.cli.Close()
	}
}

func (r *Runtime) Quiet() bool {
	return r.quiet
}

func (r *Runtime) Root() string {
	if len(r.root_path) > 0 {
		return r.root_path
	}

	dir := filepath.Dir(r.config_file)

	var err error

	if r.root_path, err = filepath.Abs(dir); err != nil {
		panic(err)
	}

	return r.root_path
}

func (r *Runtime) GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	root := r.Root()
	full_path := filepath.Join(root, path)

	return filepath.Clean(full_path)
}

func (r *Runtime) Config() *contract.Config {
	if r.config != nil {
		return r.config
	}

	var err error

	if r.config, err = getConfigFromFile(r.config_file); err != nil {
		panic(err)
	}

	return r.config
}

func (r *Runtime) Terminal() *cli.CLI {
	if r.rpc != nil {
		r.rpc.Close()
	}

	return r.getCLI(false)
}

func (r *Runtime) CLI() *cli.CLI {
	if r.rpc != nil {
		r.rpc.Close()
	}

	return r.getCLI(true)
}

func (r *Runtime) RPC() *rpc.RPC {
	if r.rpc != nil {
		return r.rpc
	}

	cli_session := r.CLI()

	if err := cli_session.SendCommand("start_rpc_session"); err != nil {
		panic(err)
	}

	if _, found, err := cli_session.ReadUntil(cli.CRLF); !found {
		panic("unable to init RPC session")
	} else if err != nil {
		panic(err)
	}

	r.rpc = rpc.New(cli_session, func() {
		if err := r.rpc.StopSession(); err != nil {
			panic(err)
		} else {
			r.rpc = nil
		}

		if _, err := cli_session.ReadUntilTerminal(); err != nil {
			panic(err)
		}
	})

	return r.rpc
}

func (ctx *Runtime) getCLI(skip_splash bool) *cli.CLI {
	if ctx.cli != nil {
		return ctx.cli
	}

	var err error
	var port string

	if len(ctx.port) > 0 {
		port = ctx.port
	} else if port, err = cli.GetFlipperPort(); err != nil {
		panic(err)
	}

	if ctx.cli, err = cli.Open(port); err != nil {
		panic(err)
	}

	if !skip_splash {
		return ctx.cli
	}

	if _, err := ctx.cli.ReadUntilTerminal(); err != nil {
		panic(err)
	}

	return ctx.cli
}
