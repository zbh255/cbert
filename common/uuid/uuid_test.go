package uuid

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestExample(t *testing.T) {
	Int64 := [8]byte{}
	Int64[0] = 128
	buf := uint64(0)
	// little endian handling
	for k := range Int64 {
		*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buf)) + uintptr(7-k))) = Int64[k]
	}
	uBuf := *(*uint64)(unsafe.Pointer(&Int64))
	slice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&Int64)),
		Len:  len(Int64),
		Cap:  len(Int64),
	}))
	fmt.Println(binary.BigEndian.Uint64(slice))
	fmt.Println(buf)
	fmt.Println(uBuf)
}

func TestUuidGenerator(t *testing.T) {
	// testing 1 million equals
	hashTable := make(map[string]struct{}, 100000)
	for i := 0; i < 100000; i++ {
		uuid, err := GetCustomUuid()
		if err != nil {
			t.Error(err)
			return
		}
		_, ok := hashTable[uuid]
		if ok {
			t.Error(errors.New("uuid generation failed"))
		}
		hashTable[uuid] = struct{}{}
	}
	// check uuid
	for k := range hashTable {
		_, err := DecodeUuid(k)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
