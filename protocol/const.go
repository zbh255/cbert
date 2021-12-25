package protocol

// define std cbert connection protocol constants

const (
	VERSION byte = 0x01
)

const (
	REPLY_SUCCESS byte = 0x00
	REPLY_FAILED  byte = 0x01
)

const (
	CMD_GET_SIGNLE_DATA byte = 0x00
	CMD_GET_MUTIL_DATA  byte = 0x01
)
