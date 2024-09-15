package rpc

import (
	"github.com/ofabel/fssdk/rpc/protobuf/flipper"
)

func (rpc *RPC) StopSession() error {
	request := &flipper.Main{
		Content: &flipper.Main_StopSession{
			StopSession: &flipper.StopSession{},
		},
	}

	_, err := rpc.sendAndReceive(request)

	return err
}
