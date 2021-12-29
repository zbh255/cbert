package utils

import (
	"io"
	"reflect"
	"unsafe"
)

const (
	BUFFER_SIZE int = 512
)

// link : https://github.com/zbh255/ss5-simple/blob/main/net/connect.go#L65
func ReadAll(reader io.Reader) ([]byte,error) {
	readN := 0
	buf := make([]byte, BUFFER_SIZE)
	n, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}
	for n == BUFFER_SIZE {
		readN += n
		buf = append(buf, make([]byte, BUFFER_SIZE)...)
		n, err = reader.Read(buf[readN : readN+BUFFER_SIZE])
		if err == io.EOF {
			n = 0
			break
		}
		if err != nil {
			return nil, err
		}
	}
	readN += n

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Len = readN
	return *(*[]byte)(unsafe.Pointer(sh)), nil
}