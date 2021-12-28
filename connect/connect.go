package connect

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/zbh255/cbert/common/uuid"
	"github.com/zbh255/cbert/ioc"
	"github.com/zbh255/cbert/protocol"
	"io/ioutil"
	"net"
	"sync"
)

// the library handler server connection

type Connection struct {
	listener net.Listener
	// uuid pool
	pool sync.Map
}

func NewConnection(listener net.Listener) *Connection {
	return &Connection{
		listener: listener,
		pool:     sync.Map{},
	}
}

func (c *Connection) Start() error {
	stdLog := ioc.GetStdLogger()
	errLog := ioc.GetErrorLogger()
	defer c.listener.Close()
	stdLog.Info("start accept")
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			stdLog.Error("client connection error: " + err.Error())
			return err
		}
		connAddr := conn.RemoteAddr()
		stdLog.Info(fmt.Sprintf("client open connection : [%s] -> %s",connAddr.Network(),connAddr.String()))
		go func() {
			err := c.handlerConnections(conn)
			if err != nil {
				errLog.Error("connect.Start:" + err.Error())
			}
		}()
	}
}

// close connection and delete map hash key
func (c *Connection) CloseConnection(conn net.Conn) error {
	c.pool.Delete(conn)
	stdLog := ioc.GetStdLogger()
	connAddr := conn.RemoteAddr()
	stdLog.Info(fmt.Sprintf("client exit connection : [%s] -> %s",connAddr.Network(),connAddr.String()))
	return conn.Close()
}

func (c *Connection) handlerConnections(conn net.Conn) error {
	defer c.CloseConnection(conn)
	if err := c.handlerHandshake(conn); err != nil {
		return err
	}
	err := c.handlerMessage(conn)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) handlerHandshake(conn net.Conn) error {
	request, err := protocol.DecodeHandshakeRequest(conn)
	if err != nil {
		return err
	}
	// check uuid
	userConfig := ioc.GetUserConfig()
	var checkOkUuid string
	for _, user := range userConfig.User {
		uuidBytes, err := uuid.DecodeUuid(user)
		if err != nil {
			return err
		}
		if bytes.Equal(uuidBytes, request.UuidV4[:32]) {
			checkOkUuid = user
			break
		}
	}
	response := new(protocol.HandshakeResponse)
	response.Version = protocol.VERSION
	response.Reserved = 0x00
	// no check ok uuid
	if checkOkUuid == "" {
		response.Reply = protocol.REPLY_AUTH_FAILED
		return errors.New("uuid authentication failed")
	}
	responseBytes := protocol.EncodeHandshakeResponse(response)
	n, err := conn.Write(responseBytes)
	if err != nil {
		return err
	}
	if n != len(responseBytes) {
		return protocol.ErrWriteBytesNoEqual
	}
	c.pool.Store(conn, checkOkUuid)
	return nil
}

func (c *Connection) handlerMessage(conn net.Conn) error {
	for {
		request, err := protocol.DecodeMessageRequest(conn)
		if err != nil {
			return err
		}
		response := new(protocol.MessageResponse)
		response.Version = protocol.VERSION
		response.Reserved = 0x00
		// read file
		uuidStr,ok := c.pool.Load(conn)
		if !ok {
			return errors.New("UUID data in pool not found")
		}
		projectConfig := ioc.GetProjectConfig()
		fileBytes,err := ioutil.ReadFile(projectConfig.Source.Dir + "/" + uuidStr.(string) + "/" + string(request.FileName))
		if err != nil {
			response.Reply = protocol.REPLY_NOT_FILE
		} else {
			response.Reply = protocol.REPLY_SUCCESS
		}
		response.FileName = request.FileName
		response.FileNameLength =request.FileNameLength
		// encode
		responseBytes := protocol.EncodeMessageResponse(response)
		n,err := conn.Write(responseBytes)
		if err != nil {
			return err
		}
		if n != len(responseBytes) {
			return protocol.ErrWriteBytesNoEqual
		}
		// Write data
		switch response.Reply {
		case protocol.REPLY_SUCCESS:
			n, err = conn.Write(fileBytes)
			if err != nil {
				return err
			}
			if n != len(fileBytes) {
				return errors.New("send to client file byte number not equal")
			}
		case protocol.REPLY_NOT_FILE,protocol.REPLY_FAILED:
			continue
		}
	}
}
