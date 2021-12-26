package protocol

import "io"

func DecodeHandshakeRequest(reader io.Reader) (*HandshakeRequest, error) {
	// 1 + 1 + 32
	buffer := make([]byte,34)
	n,err := reader.Read(buffer)
	if err != nil {
		return nil,err
	}
	if n != len(buffer) {
		return nil,ErrProtocolFormatFailed
	}
	request := new(HandshakeRequest)
	request.Version = buffer[0]
	request.Reserved = buffer[1]
	copy(request.UuidV4[:32], buffer[2:])
	return request,nil
}

func EncodeHandshakeResponse(rep *HandshakeResponse) []byte {
	buffer := make([]byte,0,3)
	buffer = append(buffer,rep.Version)
	buffer = append(buffer,rep.Reserved)
	buffer = append(buffer,rep.Reply)
	return buffer
}

func DecodeMessageRequest(reader io.Reader) (*MessageRequest, error) {
	buf1 := make([]byte,3)
	n,err := reader.Read(buf1)
	if err != nil {
		return nil, err
	}
	if n != len(buf1) {
		return nil, ErrProtocolFormatFailed
	}
	request := new(MessageRequest)
	request.Version = buf1[0]
	request.Command= buf1[1]
	request.FileNameLength = buf1[2]

	buf2 := make([]byte,request.FileNameLength)
	n,err = reader.Read(buf2)
	if err != nil {
		return nil, err
	}
	if n != len(buf2) {
		return nil, ErrProtocolFormatFailed
	}
	request.FileName = buf2
	return request, err
}
