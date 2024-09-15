package base

import (
	"strings"

	"github.com/ofabel/fssdk/contract"
)

func CleanPath(path string) string {
	path = strings.ReplaceAll(path, "\\", contract.DirSeparator)
	path = strings.TrimPrefix(path, contract.ExtStorageBasePath+contract.DirSeparator)
	path = strings.TrimRight(path, contract.DirSeparator)

	parts := strings.Split(path, contract.DirSeparator)

	path = contract.ExtStorageBasePath

	for _, part := range parts {
		part = strings.Trim(part, " ")

		if len(part) > 0 && part != contract.ThisDirectory {
			path += contract.DirSeparator
			path += part
		}
	}

	return path
}
