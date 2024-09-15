package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

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
