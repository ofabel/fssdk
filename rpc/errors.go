package rpc

import (
	"errors"

	"github.com/ofabel/fssdk/rpc/protobuf/flipper"
)

var ErrUnknown = errors.New("unknown error")
var ErrDecode = errors.New("decode error")
var ErrNotImplemented = errors.New("not implemented error")
var ErrBusy = errors.New("busy error")
var ErrContinuousCommandInterrupted = errors.New("continuous command interrupted error")
var ErrInvalidParameters = errors.New("invalid parameters error")
var ErrStorageNotReady = errors.New("storage not ready error")
var ErrStorageExist = errors.New("storage exist error")
var ErrStorageNotExist = errors.New("storage not exist error")
var ErrStorageInvalidParameter = errors.New("storage invalid parameter error")
var ErrStorageDenied = errors.New("storage denied error")
var ErrStorageInvalidName = errors.New("storage invalid name error")
var ErrStorageInternal = errors.New("storage internal error")
var ErrStorageNotImplemented = errors.New("storage not implemented error")
var ErrStorageAlreadyOpen = errors.New("storage already open error")
var ErrStorageDirNotEmpty = errors.New("storage dir not empty error")
var ErrAppCantStart = errors.New("app cant start error")
var ErrAppSystemLocked = errors.New("app system locked error")
var ErrAppNotRunning = errors.New("app not running error")
var ErrAppCmdError = errors.New("app cmd error error")
var ErrVirtualDisplayAlreadyStarted = errors.New("virtual display already started error")
var ErrVirtualDisplayNotStarted = errors.New("virtual display not started error")
var ErrGpioModeIncorrect = errors.New("gpio mode incorrect error")
var ErrGpioUnknownPinMode = errors.New("gpio unknown pin mode error")

var (
	errorCodeMapping = map[flipper.CommandStatus]error{
		1:  ErrUnknown,
		2:  ErrDecode,
		3:  ErrNotImplemented,
		4:  ErrBusy,
		14: ErrContinuousCommandInterrupted,
		15: ErrInvalidParameters,
		5:  ErrStorageNotReady,
		6:  ErrStorageExist,
		7:  ErrStorageNotExist,
		8:  ErrStorageInvalidParameter,
		9:  ErrStorageDenied,
		10: ErrStorageInvalidName,
		11: ErrStorageInternal,
		12: ErrStorageNotImplemented,
		13: ErrStorageAlreadyOpen,
		18: ErrStorageDirNotEmpty,
		16: ErrAppCantStart,
		17: ErrAppSystemLocked,
		21: ErrAppNotRunning,
		22: ErrAppCmdError,
		19: ErrVirtualDisplayAlreadyStarted,
		20: ErrVirtualDisplayNotStarted,
		58: ErrGpioModeIncorrect,
		59: ErrGpioUnknownPinMode,
	}
)
