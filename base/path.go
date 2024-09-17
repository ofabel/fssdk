package base

import (
	"strings"

	"github.com/ofabel/fssdk/contract"
)

func Flipper_GetCleanPath(parts ...string) string {
	path := ""

	for _, part := range parts {
		part = strings.ReplaceAll(part, "\\", contract.DirSeparator)
		part = strings.Trim(part, contract.DirSeparator)
		part = strings.Trim(part, " ")

		if len(part) > 0 {
			path += contract.DirSeparator + part
		}
	}

	path = strings.TrimPrefix(path, contract.ExtStorageBasePath)
	path = strings.Trim(path, contract.DirSeparator)

	segments := strings.Split(path, contract.DirSeparator)

	path = contract.ExtStorageBasePath

	for _, segment := range segments {
		segment = strings.Trim(segment, " ")

		if len(segment) > 0 && segment != contract.ThisDirectory {
			path += contract.DirSeparator
			path += segment
		}
	}

	return path
}

func Flipper_GetCleanPathWithoutStorage(parts ...string) string {
	path := Flipper_GetCleanPath(parts...)

	path = strings.TrimPrefix(path, contract.ExtStorageBasePath)
	path = strings.Trim(path, contract.DirSeparator)

	return path
}
