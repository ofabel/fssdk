package rpc

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ofabel/fssdk/base"
	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc/protobuf/flipper"
	"github.com/ofabel/fssdk/rpc/protobuf/storage"
)

type ProgressHandler func(progress float32)

const ChunkSize = 1024

var ErrNoRegularFile = errors.New("no regular file")

func (rpc *RPC) Storage_WalkFiles(root string, path string, walker contract.FileWalker) error {
	root_path := base.CleanFlipperPath(root)
	base_path := base.CleanFlipperPathWithoutStorage(path)
	clean_path := base.CleanFlipperPath(root_path + contract.DirSeparator + base_path)

	request := &flipper.Main{
		Content: &flipper.Main_StorageListRequest{
			StorageListRequest: &storage.ListRequest{
				Path: clean_path,
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
				Path: clean_path + contract.DirSeparator + file.Name,
				Dir:  clean_path,
				Rel:  base.CleanFlipperPathWithoutStorage(base_path + contract.DirSeparator + file.Name),
				Size: int64(file.Size),
			}); err != nil {
				return err
			}

			continue
		}

		if file.Type != storage.File_DIR {
			continue
		}

		err := rpc.Storage_WalkFiles(root_path, base_path+contract.DirSeparator+file.Name, walker)

		if err != nil {
			return err
		}
	}

	return nil
}

func (rpc *RPC) Storage_GetTree(path string) ([]*contract.File, error) {
	files := make([]*contract.File, 0, 32)

	err := rpc.Storage_WalkFiles(path, "", func(file *contract.File) error {
		files = append(files, file)

		return nil
	})

	return files, err
}

func (rpc *RPC) Storage_UploadFile(source string, target string, on_progress ProgressHandler) error {
	stat, err := os.Stat(source)

	if err != nil {
		return err
	}

	if !stat.Mode().IsRegular() {
		return ErrNoRegularFile
	}

	if same, _ := rpc.Storage_CheckFilesAreSame(source, target); same {
		return nil
	}

	fp, err := os.Open(source)

	if err != nil {
		return err
	}

	defer fp.Close()

	chunk := 0
	buffer := make([]byte, ChunkSize)
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

		on_progress(progress)
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

	file := response.GetStorageStatResponse().GetFile()

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

func (rpc *RPC) Storage_FolderExists(path string) (bool, error) {
	request := &flipper.Main{
		Content: &flipper.Main_StorageStatRequest{
			StorageStatRequest: &storage.StatRequest{
				Path: path,
			},
		},
	}

	response, err := rpc.sendAndReceive(request)

	if err != nil && response != nil && response.CommandStatus == flipper.CommandStatus_ERROR_STORAGE_NOT_EXIST {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return response.GetStorageStatResponse().GetFile().GetType() == storage.File_DIR, nil
}

func (rpc *RPC) Storage_CreateFolder(path string) error {
	request := &flipper.Main{
		Content: &flipper.Main_StorageMkdirRequest{
			StorageMkdirRequest: &storage.MkdirRequest{
				Path: path,
			},
		},
	}

	_, err := rpc.sendAndReceive(request)

	return err
}

func (rpc *RPC) Storage_CreateFolderRecursive(path string) error {
	clean_path := base.CleanFlipperPath(path)

	if exists, err := rpc.Storage_FolderExists(clean_path); exists || err != nil {
		return err
	}

	clean_path = base.CleanFlipperPathWithoutStorage(path)

	parts := strings.Split(clean_path, contract.DirSeparator)

	clean_path = contract.ExtStorageBasePath

	for _, part := range parts {
		clean_path += contract.DirSeparator + part

		if exists, err := rpc.Storage_FolderExists(clean_path); err != nil {
			return err
		} else if exists {
			continue
		}

		if err := rpc.Storage_CreateFolder(clean_path); err != nil {
			return err
		}
	}

	return nil
}

func (rpc *RPC) Storage_DeleteRecursive(path string) error {
	clean_path := base.CleanFlipperPath(path)

	request := &flipper.Main{
		Content: &flipper.Main_StorageDeleteRequest{
			StorageDeleteRequest: &storage.DeleteRequest{
				Path:      clean_path,
				Recursive: true,
			},
		},
	}

	_, err := rpc.sendAndReceive(request)

	return err
}
