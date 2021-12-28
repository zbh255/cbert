package test

import (
	"errors"
	"fmt"
	"github.com/zbh255/cbert/common/uuid"
	"github.com/zbh255/cbert/ioc"
	"github.com/zbh255/cbert/protocol"
	"github.com/zbh255/cbert/utils"
	"net"
	"os"
	"testing"
)

// Run the server and then run the test
// test running a path : ./build
func TestServer(t *testing.T) {
	err := os.Chdir("../build")
	if err != nil {
		panic(err)
	}
	err = ioc.RegisterProjectConfig("./project_config.json")
	if err != nil {
		t.Error(err)
		return
	}
	err = ioc.RegisterUserConfig(ioc.GetProjectConfig().Source.UserConfig)
	if err != nil {
		t.Error(err)
		return
	}
	// create source
	projectConfig := ioc.GetProjectConfig()
	userConfig := ioc.GetUserConfig()
	uuidStr := userConfig.User[0]
	fileName := "test_file0x77.json"
	filePath := fmt.Sprintf("%s/%s/%s",projectConfig.Source.Dir,uuidStr,fileName)
	file,err := os.OpenFile(filePath,os.O_CREATE | os.O_RDWR,0755)
	defer os.Remove(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	// writer test data
	_, err = file.Write([]byte("{\n  \"404\": \"没有该文件\",\n  \"200\": \"OK\"\n}"))
	if err != nil {
		t.Error(err)
		return
	}
	// connection server
	client,err := net.Dial("tcp",ioc.GetProjectConfig().Connection.Addr)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()
	handshakeRequest := new(protocol.HandshakeRequest)
	handshakeRequest.Version = protocol.VERSION
	handshakeRequest.Reserved = 0x00
	uuidSplitBytes,err := uuid.DecodeUuid(uuidStr)
	if err != nil {
		t.Error(err)
		return
	}
	copy(handshakeRequest.UuidV4[:],uuidSplitBytes)
	// encode request
	handshakeRequestBytes := protocol.EncodeHandshakeRequest(handshakeRequest)
	n, err := client.Write(handshakeRequestBytes)
	if err != nil {
		t.Error(err)
		return
	}
	if n != len(handshakeRequestBytes) {
		t.Error(protocol.ErrWriteBytesNoEqual)
		return
	}
	// read and decode handshake response
	handshakeResponse,err := protocol.DecodeHandshakeResponse(client)
	if err != nil {
		t.Error(err)
		return
	}
	if handshakeResponse.Reply != protocol.REPLY_SUCCESS {
		t.Error(errors.New("server authentication failed"))
		return
	}
	// encode message request
	messageRequest := new(protocol.MessageRequest)
	messageRequest.Version = protocol.VERSION
	messageRequest.Command = protocol.CMD_GET_SIGNLE_DATA
	messageRequest.FileNameLength = byte(len(fileName))
	messageRequest.FileName = []byte(fileName)

	messageRequestBytes := protocol.EncodeMessageRequest(messageRequest)
	_, err = client.Write(messageRequestBytes)
	if err != nil {
		t.Error(err)
		return
	}
	// read message response
	messageResponse, err := protocol.DecodeMessageResponse(client)
	if err != nil {
		t.Error(err)
		return
	}
	if messageResponse.Reply != protocol.REPLY_SUCCESS {
		t.Error(errors.New("server reply incorrect"))
		return
	}
	// read data
	fileData,err := utils.ReadAll(client)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(fileData))
}
