package base

import (
	"errors"
	"os"
)

var ErrStdinEmpty = errors.New("no data from STDIN available")

func ReadFromStdin(buffer []byte) (int, error) {
	if info, err := os.Stdin.Stat(); err != nil {
		return 0, err
	} else if mode := info.Mode(); mode&os.ModeNamedPipe == 0 {
		return 0, ErrStdinEmpty
	}

	return os.Stdin.Read(buffer)
}
