package cli

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/ofabel/fssdk/cli"
)

func readFromSerial(session *cli.CLI, screen tcell.Screen) {
	lines := make([]string, 1, 100)
	cx, cy, x, y := 0, 0, 0, 0
	line := ""

	for {
		lines[0] = line
		// update screen
		for i := range x {
			var c rune = 0

			if i < len(line) {
				c = rune(line[i])
			}

			screen.SetContent(i, y, c, nil, tcell.StyleDefault)
		}

		screen.ShowCursor(cx, cy)

		screen.Show()

		// read & handle character
		chr := readCharacter(session)

		if chr == -1 {
			return
		}

		if chr == '\r' {
			cx = 0

			continue
		}

		if chr == '\n' {
			cx = 0
			cy++

			x = 0
			y++

			line = ""

			continue
		}

		if chr == Delete {
			cx--

			line = line[:cx] + line[cx+1:]

			continue
		}

		// simple VT100 parsing, see https://vt100.net/docs/vt100-ug/chapter3.html
		if chr == Escape {
			chr = readCharacter(session)

			if chr == '[' {
				num, chr := readNumUntilRune(session)

				if chr == ';' {
					num2, chr := readNumUntilRune(session)

					if chr == 'f' {
						line = ""

						y = 0
						x = 0

						cx = num2
						cy = num
					}

					continue
				}

				// left arrow
				if chr == 'D' {
					cx--

					continue
				}

				// right arrow
				if chr == 'C' {
					cx++

					continue
				}

				// backspace / delete
				if chr == 'P' {
					line = line[:cx] + line[cx+1:]

					continue
				}

				// erase line
				if chr == 'K' && num == 2 {
					line = ""

					cx = 0

					continue
				}

				// erase page
				if chr == 'J' && num == 2 {
					line = ""

					y = 0
					x = 0

					cx = 0

					continue
				}

				// insert start
				if chr == 'h' && num == 4 {
					// NOP

					continue
				}

				// insert end
				if chr == 'l' && num == 4 {
					// NOP

					continue
				}
			}

			continue
		}

		// skip all non printable characters
		if chr < 32 || chr > 126 {
			// NOP

			continue
		}

		line = line[:cx] + string(chr) + line[cx:]

		x++
		cx++
	}
}

func readNumUntilRune(session *cli.CLI) (int, rune) {
	var num int = 0

	for {
		next := readCharacter(session)

		if next >= '0' && next <= '9' {
			num *= 10
			num += (int(next) - '0')
		} else {
			return num, next
		}
	}
}

func readCharacter(session *cli.CLI) rune {
	buffer := make([]byte, 1)

	for {
		if n, err := session.Read(buffer); err != nil {
			panic(err)
		} else if n > 0 {
			return rune(buffer[0])
		} else {
			time.Sleep(time.Microsecond)
		}

		if exit.Load() {
			return -1
		}
	}
}
