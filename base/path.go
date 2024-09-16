package base

import (
	"strings"

	"github.com/ofabel/fssdk/contract"
)

func CleanFlipperPath(path string) string {
	clean_path := strings.ReplaceAll(path, "\\", contract.DirSeparator)

	clean_path = strings.TrimPrefix(clean_path, contract.ExtStorageBasePath)
	clean_path = strings.Trim(clean_path, contract.DirSeparator)

	parts := strings.Split(clean_path, contract.DirSeparator)

	clean_path = contract.ExtStorageBasePath

	for _, part := range parts {
		part = strings.Trim(part, " ")

		if len(part) > 0 && part != contract.ThisDirectory {
			clean_path += contract.DirSeparator
			clean_path += part
		}
	}

	return clean_path
}

func CleanFlipperPathWithoutStorage(path string) string {
	clean_path := CleanFlipperPath(path)

	clean_path = strings.TrimPrefix(clean_path, contract.ExtStorageBasePath)
	clean_path = strings.Trim(clean_path, contract.DirSeparator)

	return clean_path
}
