package cli

import (
	"github.com/gdamore/tcell"
)

type terminal struct {
	w      int
	h      int
	cx     int
	cy     int
	line   string
	lines  []string
	screen tcell.Screen
	redraw bool
}

func NewTerminal(screen tcell.Screen) *terminal {
	w, h := screen.Size()
	lines := make([]string, 1, h)
	lines[0] = "-- Press <ESC> to finish this terminal session. --"

	return &terminal{
		w:      w,
		h:      h,
		cx:     0,
		cy:     1,
		line:   "",
		lines:  lines,
		screen: screen,
		redraw: false,
	}
}

func (t *terminal) Render() {
	if t.redraw {
		for y, line := range t.lines {
			t.RenderLine(line, y)
		}
	}

	t.RenderLine(t.line, t.cy)

	t.screen.ShowCursor(t.cx, t.cy)

	t.screen.Show()

	t.redraw = false
}

func (t *terminal) RenderLine(line string, y int) {
	for x := range t.w {
		var c rune = 0

		if x < len(line) {
			c = rune(line[x])
		}

		t.screen.SetContent(x, y, c, nil, tcell.StyleDefault)
	}
}

func (t *terminal) Resize() {
	t.w, t.h = t.screen.Size()

	t.redraw = true

	t.Render()
}

func (t *terminal) CarriageReturn() {
	t.cx = 0
}

func (t *terminal) LineFeed() {
	t.cx = 0
	t.cy++

	t.lines = append(t.lines, t.line)

	if len(t.lines) >= t.h {
		t.lines = t.lines[1:]

		t.redraw = true
		t.cy--
	}

	t.line = ""
}

func (t *terminal) Backspace() {
	t.cx--

	t.line = t.line[:t.cx] + t.line[t.cx+1:]
}

func (t *terminal) SetCursor(x int, y int) {
	if y < len(t.lines) {
		t.line = t.lines[y]
	}

	t.cx = x
	t.cy = y
}

func (t *terminal) MoveLeft() {
	t.cx--
}

func (t *terminal) MoveRight() {
	t.cx++
}

func (t *terminal) MoveUp() {
	t.cy--
}

func (t *terminal) MoveDown() {
	t.cy++
}

func (t *terminal) EraseCharacter() {
	t.line = t.line[:t.cx] + t.line[t.cx+1:]
}

func (t *terminal) EraseLine() {
	t.line = ""

	t.cx = 0
}

func (t *terminal) ErasePage() {
	t.line = ""
	t.lines = make([]string, 0, t.h)

	t.cx = 0
	t.cy = 0

	t.screen.Clear()
}

func (t *terminal) Insert(chr rune) {
	t.line = t.line[:t.cx] + string(chr) + t.line[t.cx:]

	t.cx++
}
