package cli

import (
	"strings"
	"sync/atomic"

	"github.com/gdamore/tcell"
	"github.com/ofabel/fssdk/app"
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

	terminal := NewTerminal(screen)

	go readFromSerial(session, terminal)

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
		case *tcell.EventResize:
			terminal.Resize()
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

	if _, err := session.ReadUntilTerminal(); err != nil {
		panic(err)
	}

	session.SendCommand(args.Command)

	if raw_output, err := session.ReadUntilTerminal(); err != nil {
		panic(err)
	} else {
		n := len(args.Command)
		output := string(raw_output[n:])
		output = strings.Trim(output, "\r\n")

		runtime.Printf("%s\n", output)
	}
}
