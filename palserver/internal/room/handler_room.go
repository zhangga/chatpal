package room

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/zhangga/chatpal/palserver/internal/msg"
	"time"
)

type CreateRoomReq struct {
	Name string `json:"Name,omitempty"`
	Desc string `json:"Desc,omitempty"`
}
type CreateRoomResp struct {
	Code   int   `json:"code,omitempty"`
	RoomId int64 `json:"room_id,omitempty"`
}

// handleCreateRoom 创建房间
func handleCreateRoom(c *conn, msgType int, content []byte) (msg.Message, error) {
	var req CreateRoomReq
	if err := json.Unmarshal(content, &req); err != nil {
		c.logger.Errorf("error while unmarshaling create room request, err: %v", err)
		return CreateRoomResp{Code: 1}, nil
	}
	if len(req.Name) == 0 {
		c.logger.Infof("room Name is empty")
		return CreateRoomResp{Code: 2}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	id, err := c.redisClient.Incr(ctx, rkeyRoomId).Result()
	if err != nil {
		c.logger.Errorf("error while generating room Id, err: %v", err)
		return CreateRoomResp{Code: 3}, nil
	}

	room := NewRoom(id, req.Name, req.Desc)
	room.OwnerOpenId = c.user.OpenId
	if err = room.Save(c.redisClient); err != nil {
		c.logger.Errorf("error while saving room info, err: %v", err)
		return CreateRoomResp{Code: 4}, nil
	}
	return CreateRoomResp{Code: 0, RoomId: room.Id}, nil
}

// handleListRoom 房间列表
func handleListRoom(c *conn, msgType int, content []byte) (msg.Message, error) {
	return nil, nil
}
