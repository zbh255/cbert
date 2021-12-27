package protocol

type HandshakeRequest struct {
	Version  byte
	Reserved byte
	UuidV4   [32]byte
}

type HandshakeResponse struct {
	Version  byte
	Reserved byte
	Reply    byte
}

type MessageRequest struct {
	Version        byte
	Command        byte
	FileNameLength byte
	FileName       []byte
}

type MessageResponse struct {
	Version        byte
	Reserved       byte
	Reply          byte
	FileNameLength byte
	FileName       []byte
}
