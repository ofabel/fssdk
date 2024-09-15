package flipper

import (
	"errors"
	"strings"

	"github.com/albenik/go-serial/v2/enumerator"
)

var ErrNoFlipperFound = errors.New("no flipper device found")

func GetFlipperPort() (string, error) {
	ports, err := enumerator.GetDetailedPortsList()

	if err != nil {
		return "", err
	}
	if len(ports) == 0 {
		return "", ErrNoFlipperFound
	}
	for _, port := range ports {
		if strings.HasPrefix(port.SerialNumber, "flip_") {
			return port.Name, nil
		}
	}

	return "", ErrNoFlipperFound
}
