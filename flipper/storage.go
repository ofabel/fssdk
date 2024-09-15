package flipper

import (
	"github.com/ofabel/fssdk/flipper/rpc/flipper"
	"github.com/ofabel/fssdk/flipper/rpc/storage"
)

type File struct {
	Name string
	Path string
	Size uint32
}

type FileWalker func(file *File)

func (f0 *Flipper) WalkStorageFiles(path string, walker FileWalker) error {
	request := &flipper.Main{
		Content: &flipper.Main_StorageListRequest{
			StorageListRequest: &storage.ListRequest{
				Path: path,
			},
		},
	}

	seq, err := f0.send(request)

	if err != nil {
		return err
	}

	collected_files := make([]*storage.File, 0, 32)

	for {
		response, err := f0.readAnswer(seq)

		if err != nil {
			return err
		}

		var files = response.GetStorageListResponse().GetFile()

		collected_files = append(collected_files, files...)

		if !response.HasNext {
			break
		}
	}

	for _, file := range collected_files {
		if file.Type == storage.File_FILE {
			walker(&File{
				Name: file.Name,
				Path: path + "/" + file.Name,
				Size: file.Size,
			})

			continue
		}

		if file.Type != storage.File_DIR {
			continue
		}

		err := f0.WalkStorageFiles(path+"/"+file.Name, walker)

		if err != nil {
			return err
		}
	}

	return nil
}
