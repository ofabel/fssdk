//go:build darwin

package cli

func GetFlipperPort() (string, error) {
	return "", ErrNoFlipperFound
}
