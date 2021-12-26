package uuid

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"time"
	"unsafe"
)

var (
	/*
		Available MAC addr
		If there is no network hardware, the value is 0
	*/
	hardwareAddr []byte
)


// GetCustomUuid UUID format
// 64 bit mac addr XOR random seed | 64 bit Random number
// run in a use little endian byte order computer
func GetCustomUuid() (string, error) {

	macAddr64 := binary.LittleEndian.Uint64(hardwareAddr)
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

func init() {
	interfaces,err := net.Interfaces()
	hardwareAddr = make([]byte,8)
	// no network interfaces
	if len(interfaces) == 0 || err != nil {
		return
	}
	for _,v := range interfaces {
		if v.HardwareAddr != nil {
			copy(hardwareAddr,v.HardwareAddr)
			return
		}
	}
}