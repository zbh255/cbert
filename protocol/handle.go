package protocol

import "io"

func DecodeHandshakeRequest(reader io.Reader) (*HandshakeRequest, error) {
	// 1 + 1 + 32
	buffer := make([]byte, 34)
	n, err := reader.Read(buffer)
	if err != nil {
		return nil, err
	}
	if n != len(buffer) {
		return nil, ErrProtocolFormatFailed
	}
	request := new(HandshakeRequest)
	request.Version = buffer[0]
	request.Reserved = buffer[1]
	copy(request.UuidV4[:32], buffer[2:])
	return request, nil
}

// suport client
func EncodeHandshakeRequest(request *HandshakeRequest) []byte {
	if request == nil {
		return nil
	}
	buffer := make([]byte,0,34)
	buffer = append(buffer,request.Version)
	buffer = append(buffer,request.Reserved)
	buffer = append(buffer,request.UuidV4[:]...)
	return buffer
}

func EncodeHandshakeResponse(rep *HandshakeResponse) []byte {
	buffer := make([]byte, 0, 3)
	buffer = append(buffer, rep.Version)
	buffer = append(buffer, rep.Reserved)
	buffer = append(buffer, rep.Reply)
	return buffer
}

func DecodeHandshakeResponse(reader io.Reader) (*HandshakeResponse, error) {
	buffer := make([]byte,3)
	n,err := reader.Read(buffer)
	if err != nil {
		return nil, err
	}
	if n != len(buffer) {
		return nil, ErrReaderBytesNoEqual
	}
	return &HandshakeResponse{
		Version:  buffer[0],
		Reserved: buffer[1],
		Reply:    buffer[2],
	}, err
}

func DecodeMessageRequest(reader io.Reader) (*MessageRequest, error) {
	buf1 := make([]byte, 3)
	n, err := reader.Read(buf1)
	if err != nil {
		return nil, err
	}
	if n != len(buf1) {
		return nil, ErrProtocolFormatFailed
	}
	request := new(MessageRequest)
	request.Version = buf1[0]
	request.Command = buf1[1]
	request.FileNameLength = buf1[2]

	buf2 := make([]byte, request.FileNameLength)
	n, err = reader.Read(buf2)
	if err != nil {
		return nil, err
	}
	if n != len(buf2) {
		return nil, ErrProtocolFormatFailed
	}
	request.FileName = buf2
	return request, err
}

func EncodeMessageRequest(request *MessageRequest) []byte {
	buf := []byte{
		request.Version,
		request.Command,
		request.FileNameLength,
	}
	buf = append(buf, request.FileName...)
	return buf
}

func EncodeMessageResponse(response *MessageResponse) []byte {
	buf := make([]byte, 0, 32)
	buf = append(buf, response.Version)
	buf = append(buf, response.Reserved)
	buf = append(buf, response.Reply)
	buf = append(buf, append([]byte{response.FileNameLength}, response.FileName...)...)
	return buf
}

func DecodeMessageResponse(reader io.Reader) (*MessageResponse, error) {
	buffer := make([]byte,4)
	n, err := reader.Read(buffer)
	if err != nil {
		return nil,err
	}
	if n != len(buffer) {
		return nil,ErrReaderBytesNoEqual
	}
	rep := &MessageResponse{
		Version:        buffer[0],
		Reserved:       buffer[1],
		Reply:          buffer[2],
		FileNameLength: buffer[3],
		FileName:       nil,
	}
	rep.FileName = make([]byte,rep.FileNameLength)
	n, err = reader.Read(rep.FileName)
	if err != nil {
		return nil, err
	}
	if n != int(rep.FileNameLength) {
		return nil,ErrReaderBytesNoEqual
	}
	return rep,nil
}