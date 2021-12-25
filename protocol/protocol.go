package protocol

type HandshakeRequest struct {
	Version  byte
	Reserved byte
	UuidV4   [16]byte
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
