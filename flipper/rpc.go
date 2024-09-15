package flipper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"syscall"

	"github.com/ofabel/fssdk/flipper/rpc/flipper"

	"google.golang.org/protobuf/proto"
)

var ErrTooManyBytes = errors.New("too many bytes when decoding varint")
var ErrTooLittleBytesWritten = errors.New("too little bytes written")

func (f0 *Flipper) getNextSeq() uint32 {
	f0.seq++

	return f0.seq
}

func (f0 *Flipper) send(request *flipper.Main) (uint32, error) {
	if request.CommandId == 0 {
		request.CommandId = f0.getNextSeq()
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

	for {
		n, err := f0.port.Write(buffer)

		// ignore EINTR
		if errors.Is(err, syscall.EINTR) {
			break
		}

		// try again on EAGAIN
		if errors.Is(err, syscall.EAGAIN) {
			continue
		}

		if err != nil {
			return 0, err
		}

		if n != len(buffer) {
			return 0, ErrTooLittleBytesWritten
		}

		break
	}

	return request.CommandId, nil
}

func (f0 *Flipper) sendAndReceive(request *flipper.Main) (*flipper.Main, error) {
	seq, err := f0.send(request)

	if err != nil {
		return nil, err
	}

	return f0.readAnswer(seq)
}

func (f0 *Flipper) readAnswer(seq uint32) (*flipper.Main, error) {
	for {
		data, err := f0.readAny()

		if err != nil {
			return nil, err
		}

		if data.CommandId == seq {
			return data, nil
		}
	}
}

func (f0 *Flipper) readAny() (*flipper.Main, error) {
	size, err := f0.readVariant32()

	if err != nil {
		return nil, err
	}

	raw_data := make([]byte, size)

	_, err = f0.port.Read(raw_data)

	if err != nil {
		return nil, err
	}

	data := &flipper.Main{}

	err = proto.Unmarshal(raw_data, data)

	if err != nil {
		return nil, err
	}

	if data.CommandStatus != flipper.CommandStatus_OK {
		return nil, fmt.Errorf("%s", data.CommandStatus)
	}

	return data, err
}

func (f0 *Flipper) readVariant32() (uint32, error) {
	const MASK = (1 << 32) - 1

	var result = uint32(0)
	var shift = uint32(0)

	var buffer = make([]byte, 1)
	var raw_data = make([]byte, 4)

	for {
		n, err := f0.port.Read(buffer)

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
