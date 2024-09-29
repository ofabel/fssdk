package cli

import (
	"time"

	"github.com/ofabel/fssdk/cli"
)

func readFromSerial(session *cli.CLI, t *terminal) {
	for {
		// update screen
		t.Render()

		// read & handle character
		chr := readCharacter(session)

		if exit.Load() {
			return
		}

		if chr == '\r' {
			t.CarriageReturn()

			continue
		}

		if chr == '\n' {
			t.LineFeed()

			continue
		}

		if chr == Delete {
			t.Backspace()

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
						t.SetCursor(num2, num)
					}

					continue
				}

				// left arrow
				if chr == 'D' {
					t.MoveLeft()

					continue
				}

				// right arrow
				if chr == 'C' {
					t.MoveRight()

					continue
				}

				// backspace / delete
				if chr == 'P' {
					t.EraseCharacter()

					continue
				}

				// erase line
				if chr == 'K' && num == 2 {
					t.EraseLine()

					continue
				}

				// erase page
				if chr == 'J' && num == 2 {
					t.ErasePage()

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

		t.Insert(chr)
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
