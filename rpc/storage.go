package rpc

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc/protobuf/flipper"
	"github.com/ofabel/fssdk/rpc/protobuf/storage"
)

type ProgressHandler func(progress float32) error

const CHUNK_SIZE = 1024

var ErrNoRegularFile = errors.New("no regular file")

func (rpc *RPC) Storage_WalkFiles(path string, walker contract.FileWalker) error {
	path = strings.TrimRight(path, "/")

	request := &flipper.Main{
		Content: &flipper.Main_StorageListRequest{
			StorageListRequest: &storage.ListRequest{
				Path: path,
			},
		},
	}

	seq, err := rpc.send(request)

	if err != nil {
		return err
	}

	collected_files := make([]*storage.File, 0, 32)

	for {
		response, err := rpc.readAnswer(seq)

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
			if err := walker(&contract.File{
				Name: file.Name,
				Path: path + "/" + file.Name,
				Size: int64(file.Size),
			}); err != nil {
				return err
			}

			continue
		}

		if file.Type != storage.File_DIR {
			continue
		}

		err := rpc.Storage_WalkFiles(path+"/"+file.Name, walker)

		if err != nil {
			return err
		}
	}

	return nil
}

func (rpc *RPC) Storage_GetTree(path string) ([]*contract.File, error) {
	files := make([]*contract.File, 0, 32)

	err := rpc.Storage_WalkFiles(path, func(file *contract.File) error {
		files = append(files, file)

		return nil
	})

	return files, err
}

func (rpc *RPC) Storage_UploadFile(source string, target string, onProgress ProgressHandler) error {
	stat, err := os.Stat(source)

	if err != nil {
		return err
	}

	if !stat.Mode().IsRegular() {
		return ErrNoRegularFile
	}

	if same, err := rpc.Storage_CheckFilesAreSame(source, target); same || err != nil {
		return err
	}

	fp, err := os.Open(source)

	if err != nil {
		return err
	}

	defer fp.Close()

	chunk := 0
	buffer := make([]byte, CHUNK_SIZE)
	size := stat.Size()
	written := int64(0)
	progress := float32(0)

	file := &storage.File{}
	request := &flipper.Main{
		CommandId: rpc.getNextSeq(),
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

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		written += int64(chunk)

		file.Data = buffer[:chunk]
		request.HasNext = written < size

		_, err = rpc.send(request)

		if err != nil {
			return err
		}

		progress = float32(written) / float32(size)

		if err := onProgress(progress); err != nil {
			return err
		}
	}

	_, err = rpc.readAnswer(request.CommandId)

	return err
}

func (rpc *RPC) Storage_GetChecksum(path string) (string, error) {
	request := &flipper.Main{
		Content: &flipper.Main_StorageMd5SumRequest{
			StorageMd5SumRequest: &storage.Md5SumRequest{
				Path: path,
			},
		},
	}

	response, err := rpc.sendAndReceive(request)

	if err != nil {
		return "", err
	}

	return response.GetStorageMd5SumResponse().Md5Sum, nil
}

func (rpc *RPC) Storage_GetFileSize(path string) (int64, error) {
	request := &flipper.Main{
		Content: &flipper.Main_StorageStatRequest{
			StorageStatRequest: &storage.StatRequest{
				Path: path,
			},
		},
	}

	response, err := rpc.sendAndReceive(request)

	if err != nil {
		return 0, err
	}

	file := response.GetStorageStatResponse().File

	if file.Type == storage.File_DIR {
		return 0, ErrNoRegularFile
	}

	return int64(file.Size), nil
}

func (rpc *RPC) Storage_CheckFilesHaveSameSize(source string, target string) (bool, error) {
	stat, err := os.Stat(source)

	if err != nil {
		return false, err
	}

	if !stat.Mode().IsRegular() {
		return false, ErrNoRegularFile
	}

	source_size := stat.Size()

	target_size, err := rpc.Storage_GetFileSize(target)

	if err != nil {
		return false, err
	}

	return source_size == target_size, nil
}

func (rpc *RPC) Storage_CheckFilesHaveSameHash(source string, target string) (bool, error) {
	stat, err := os.Stat(source)

	if err != nil {
		return false, err
	}

	if !stat.Mode().IsRegular() {
		return false, ErrNoRegularFile
	}

	source_checksum, err := getLocalFileChecksum(source)

	if err != nil {
		return false, err
	}

	target_checksum, err := rpc.Storage_GetChecksum(target)

	if err != nil {
		return false, err
	}

	return source_checksum == target_checksum, nil
}

func (rpc *RPC) Storage_CheckFilesAreSame(source string, target string) (bool, error) {
	stat, err := os.Stat(source)

	if err != nil {
		return false, err
	}

	if !stat.Mode().IsRegular() {
		return false, ErrNoRegularFile
	}

	if stat.Size() > 1024*512 {
		return rpc.Storage_CheckFilesHaveSameSize(source, target)
	} else {
		return rpc.Storage_CheckFilesHaveSameHash(source, target)
	}
}
