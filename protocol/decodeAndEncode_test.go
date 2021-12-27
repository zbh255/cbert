package protocol

import (
	"bytes"
	"errors"
	"github.com/zbh255/cbert/common/uuid"
	"testing"
)

func TestHandshake(t *testing.T) {
	uuidStr, err := uuid.GetCustomUuid()
	if err != nil {
		t.Error(err)
		return
	}
	uuidBytes, err := uuid.DecodeUuid(uuidStr)
	if err != nil {
		t.Error(err)
		return
	}
	requestBytes := append([]byte{
		VERSION,
		0x00,
	}, uuidBytes...)
	reader := bytes.NewReader(requestBytes)
	request, err := DecodeHandshakeRequest(reader)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(request.UuidV4[:32], uuidBytes) {
		t.Error(errors.New("the uuid not equal"))
	}
	t.Log(request)
	// encode response
	response := new(HandshakeResponse)
	response.Version = VERSION
	response.Reserved = 0x00
	response.Reply = REPLY_SUCCESS
	EncodeHandshakeResponse(response)
}

func TestMessage(t *testing.T) {
	fileName := []byte("/ssm-studio/exceptions.json")
	requestBytes := append([]byte{
		VERSION,
		CMD_GET_SIGNLE_DATA,
		byte(len(fileName)),
	}, fileName...)
	reader := bytes.NewReader(requestBytes)
	request, err := DecodeMessageRequest(reader)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(request)
	response := new(MessageResponse)
	response.Version = VERSION
	response.Reserved = 0x00
	response.Reply = REPLY_NOT_FILE
	response.FileNameLength = request.FileNameLength
	response.FileName = request.FileName
	responseBytes := []byte{
		VERSION,
		0x00,
		REPLY_NOT_FILE,
		27,
		47, 115, 115, 109, 45, 115, 116, 117, 100, 105, 111, 47, 101, 120, 99, 101, 112, 116, 105, 111, 110,
		115, 46, 106, 115, 111, 110,
	}
	if !bytes.Equal(responseBytes, EncodeMessageResponse(response)) {
		t.Error(errors.New("message response encode is not equal"))
		return
	}
}
