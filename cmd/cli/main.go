package cli

import (
	"sync/atomic"

	"github.com/gdamore/tcell"
	"github.com/ofabel/fssdk/app"
	"github.com/ofabel/fssdk/cli"
)

type Args struct {
	Command string `arg:"-C,--command" help:"Execute a single command."`
}

var exit = atomic.Bool{}

const Escape = 27
const Delete = 127

func Main(runtime *app.Runtime, args *Args) {
	if len(args.Command) > 0 {
		execSingleCommand(runtime, args)

		return
	}

	exit.Store(false)

	session := runtime.Terminal()

	defer session.Close()

	screen, err := tcell.NewScreen()

	if err != nil {
		panic(err)
	}

	if err := screen.Init(); err != nil {
		panic(err)
	}

	defer screen.Fini()

	go readFromSerial(session, screen)

	for {
		screen.Show()

		event := screen.PollEvent()

		switch event := event.(type) {
		case *tcell.EventKey:
			if out := handleKeyEvent(event); out != nil {
				session.Write(out)
			} else {
				exit.Store(true)
			}
		}

		if exit.Load() {
			break
		}
	}
}

func handleKeyEvent(event *tcell.EventKey) []byte {
	key := event.Key()

	if key == tcell.KeyEscape {
		return nil
	}

	if key == tcell.KeyCtrlC {
		return []byte{3}
	}

	if key == tcell.KeyCtrlD {
		return []byte{4}
	}

	if key == tcell.KeyTab {
		return []byte{9}
	}

	if key == tcell.KeyEnter {
		return []byte{'\r', '\n'}
	}

	if key == tcell.KeyBackspace {
		return []byte{8}
	}

	if key == tcell.KeyDelete {
		return []byte{Escape, '[', 'C', 8} // simle hack to support DELETE key
	}

	if key == tcell.KeyUp {
		return []byte{Escape, '[', 'A'}
	}

	if key == tcell.KeyDown {
		return []byte{Escape, '[', 'B'}
	}

	if key == tcell.KeyRight {
		return []byte{Escape, '[', 'C'}
	}

	if key == tcell.KeyLeft {
		return []byte{Escape, '[', 'D'}
	}

	return []byte{byte(event.Rune())}
}

func execSingleCommand(runtime *app.Runtime, args *Args) {
	session := runtime.Terminal()

	defer session.Close()

	if out, err := session.ReadUntilTerminal(); err != nil {
		panic(err)
	} else {
		runtime.Printf("%s%s", out, cli.TerminalDelimiter)
	}

	session.SendCommand(args.Command)

	if out, err := session.ReadUntilTerminal(); err != nil {
		panic(err)
	} else {
		runtime.Printf("%s\n", out)
	}
}
