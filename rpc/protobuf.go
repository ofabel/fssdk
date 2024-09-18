package rpc

import (
	"encoding/binary"
	"errors"

	"github.com/ofabel/fssdk/contract"
	"github.com/ofabel/fssdk/rpc/protobuf/flipper"

	"google.golang.org/protobuf/proto"
)

var ErrTooManyBytes = errors.New("too many bytes when decoding varint")
var ErrTooLittleBytesWritten = errors.New("too little bytes written")

type RPC struct {
	port     contract.IO
	seq      uint32
	on_close func()
}

func New(port contract.IO, on_close func()) *RPC {
	return &RPC{
		port:     port,
		seq:      0,
		on_close: on_close,
	}
}

func (rpc *RPC) Close() {
	rpc.on_close()

	rpc.port = nil
	rpc.seq = 0
}

func (rpc *RPC) getNextSeq() uint32 {
	rpc.seq++

	return rpc.seq
}

func (rpc *RPC) send(request *flipper.Main) (uint32, error) {
	if request.CommandId == 0 {
		request.CommandId = rpc.getNextSeq()
	}

	if request.Content == nil {
		request.Content = &flipper.Main_Empty{
			Empty: &flipper.Empty{},
		}
	}

	request.CommandStatus = flipper.CommandStatus_OK

	raw_request, err := proto.Marshal(request)

	if err != nil {
		return 0, err
	}

	var size = len(raw_request)
	var raw_size = uint64(size)
	var buffer = make([]byte, 0, binary.MaxVarintLen64)

	buffer = binary.AppendUvarint(buffer, raw_size)

	buffer = append(buffer, raw_request...)

	n, err := rpc.port.Write(buffer)

	if err != nil {
		return 0, err
	}

	if n != len(buffer) {
		return 0, ErrTooLittleBytesWritten
	}

	return request.CommandId, nil
}

func (rpc *RPC) sendAndReceive(request *flipper.Main) (*flipper.Main, error) {
	seq, err := rpc.send(request)

	if err != nil {
		return nil, err
	}

	return rpc.readAnswer(seq)
}

func (rpc *RPC) readAnswer(seq uint32) (*flipper.Main, error) {
	for {
		response, err := rpc.readAny()

		if err != nil {
			return response, err
		}

		if response.CommandId == seq {
			return response, nil
		}
	}
}

func (rpc *RPC) readAny() (*flipper.Main, error) {
	size, err := rpc.readVariant32()

	if err != nil {
		return nil, err
	}

	raw_data := make([]byte, size)

	_, err = rpc.port.Read(raw_data)

	if err != nil {
		return nil, err
	}

	response := &flipper.Main{}

	err = proto.Unmarshal(raw_data, response)

	if err != nil {
		return response, err
	}

	if response.CommandStatus != flipper.CommandStatus_OK {
		if err, ok := errorCodeMapping[response.CommandStatus]; ok {
			return response, err
		} else {
			return response, ErrUnknown
		}
	}

	return response, err
}

func (rpc *RPC) readVariant32() (uint32, error) {
	const MASK = (1 << 32) - 1

	var result = uint32(0)
	var shift = uint32(0)

	var buffer = make([]byte, 1)
	var raw_data = make([]byte, 4)

	for {
		n, err := rpc.port.Read(buffer)

		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, err
		} else {
			raw_data[0] = buffer[0]
		}

		var data = binary.LittleEndian.Uint32(raw_data)

		result |= (data & 0x7F) << shift

		if data&0x80 == 0 {
			result &= MASK

			return result, nil
		}

		shift += 7

		if shift >= 64 {
			return 0, ErrTooManyBytes
		}
	}
}
