package flipper

import (
	"errors"
	"os"
	"strings"

	"github.com/ofabel/fssdk/flipper/rpc/flipper"
	"github.com/ofabel/fssdk/flipper/rpc/storage"
)

type File struct {
	Name string
	Path string
	Size uint32
}

type FileWalker func(file *File)
type ProgressHandler func(progress float32)

const CHUNK_SIZE = 1024

var ErrNoRegularFile = errors.New("no regular file")

func (f0 *Flipper) WalkStorageFiles(path string, walker FileWalker) error {
	path = strings.TrimRight(path, "/")

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

func (f0 *Flipper) GetTree(path string) ([]*File, error) {
	files := make([]*File, 0, 32)

	err := f0.WalkStorageFiles(path, func(file *File) {
		files = append(files, file)
	})

	return files, err
}

func (f0 *Flipper) UploadFile(source string, target string, onProgress ProgressHandler) error {
	stat, err := os.Stat(source)

	if err != nil {
		return err
	}

	if !stat.Mode().IsRegular() {
		return ErrNoRegularFile
	}

	fp, err := os.Open(source)

	if err != nil {
		return err
	}

	chunk := 0
	buffer := make([]byte, CHUNK_SIZE)
	size := stat.Size()
	written := int64(0)
	progress := float32(0)

	file := &storage.File{}
	request := &flipper.Main{
		CommandId: f0.getNextSeq(),
		HasNext:   false,
		Content: &flipper.Main_StorageWriteRequest{
			StorageWriteRequest: &storage.WriteRequest{
				Path: target,
				File: file,
			},
		},
	}

	for {
		chunk, err = fp.Read(buffer)

		if err != nil {
			break
		}

		if chunk == 0 {
			break
		}

		written += int64(chunk)

		file.Data = buffer[:chunk]
		request.HasNext = written < size

		_, err = f0.send(request)

		if err != nil {
			break
		}

		progress = float32(written) / float32(size)

		onProgress(progress)
	}

	fp.Close()

	if err != nil {
		return err
	}

	_, err = f0.readAnswer(request.CommandId)

	return err
}
