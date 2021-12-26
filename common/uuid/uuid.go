package uuid

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"time"
	"unsafe"
)

// GetCustomUuid UUID format
// 64 bit mac addr XOR random seed | 64 bit Random number
// run in a use little endian byte order computer
func GetCustomUuid() (string, error) {
	interfaces,err := net.Interfaces()
	if err != nil {
		return "", err
	}
	// no network interfaces
	if len(interfaces) == 0 {
		return "", errors.New("no network interfaces")
	}
	// padding
	if n := len(interfaces[0].HardwareAddr); n < 8 {
		interfaces[0].HardwareAddr = append(interfaces[0].HardwareAddr,make([]byte,8 - n)...)
	}
	macAddr64 := binary.LittleEndian.Uint64(interfaces[0].HardwareAddr)
	seed := time.Now().UnixNano()
	one64 := macAddr64 ^ uint64(seed)
	// handle two 64
	rand.Seed(seed)
	two64 := rand.Uint64()
	buffer := [16]byte{}
	slice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&buffer)),
		Len:  8,
		Cap:  8,
	}))
	binary.LittleEndian.PutUint64(slice,one64)
	slice = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&buffer)) + uintptr(len(slice)),
		Len:  8,
		Cap:  8,
	}))
	binary.LittleEndian.PutUint64(slice,two64)
	str := hex.EncodeToString(buffer[:])
	return fmt.Sprintf("%s-%s-%s-%s-%s",str[:8],str[8:12],str[12:16],str[16:20],str[20:32]),nil
}
