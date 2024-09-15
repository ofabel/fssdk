package flipper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
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

func getLocalFileChecksum(path string) (string, error) {
	stat, err := os.Stat(path)

	if err != nil {
		return "", err
	}

	if !stat.Mode().IsRegular() {
		return "", ErrNoRegularFile
	}

	fp, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer fp.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, fp); err != nil {
		return "", err
	}

	checksum := hash.Sum(nil)

	return hex.EncodeToString(checksum), nil
}
