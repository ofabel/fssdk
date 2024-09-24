//go:build !darwin

package cli

import (
	"strings"

	"github.com/albenik/go-serial/v2/enumerator"
)

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
