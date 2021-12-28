package protocol

import "errors"

var (
	ErrProtocolFormatFailed = errors.New("error: protocol format failed")
	ErrWriteBytesNoEqual    = errors.New("writer to connection byte number not equal")
	ErrReaderBytesNoEqual = errors.New("from connection read byte number not equal")
)
