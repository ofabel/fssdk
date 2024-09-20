package sync

import (
	"slices"
	"strings"

	"github.com/ofabel/fssdk/app"
	"golang.org/x/exp/maps"
)

func ListFilesFromSyncMap(runtime *app.Runtime, sync_map SyncMap) {
	keys := maps.Keys(sync_map)

	slices.SortStableFunc(keys, compareFiles)

	for _, key := range keys {
		if file, found := sync_map[key]; found {
			runtime.Printf("[%s] %s\n", file.Status, key)
		}
	}
}

func compareFiles(a string, b string) int {
	as := strings.Split(a, "/")
	bs := strings.Split(b, "/")

	size := len(as)

	if size > len(bs) {
		size = len(bs)
	}

	for i := range size {
		if as[i] != bs[i] {
			// two files in the same folder
			if len(as) == len(bs) && size == i+1 {
				return strings.Compare(as[i], bs[i])
			}

			// as[i] is a file
			if len(as) == i+1 {
				return 1
			}

			// bs[i] is a file
			if len(bs) == i+1 {
				return -1
			}

			// compare two arbitrary folders
			return strings.Compare(as[i], bs[i])
		}
	}

	return 0
}
