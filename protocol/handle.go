package protocol

import "io"

func DecodeHandshakeRequest(io io.Reader) (*HandshakeRequest, error) {

}

func EncodeHandshakeResponse(rep *HandshakeResponse) ([]byte, error) {}

func DecodeMessageRequest(io io.Reader) (*MessageRequest, error) {}
