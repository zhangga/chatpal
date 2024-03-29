package room

import (
	"github.com/zhangga/chatpal/palserver/internal/msg"
)

type Handler func(c *conn, msgType int, content []byte) (msg.Message, error)

var msgHandlers = map[msg.Id]Handler{
	msg.Login:      handleLogin,
	msg.CreateRoom: handleCreateRoom,
	msg.ListRoom:   handleListRoom,
}
