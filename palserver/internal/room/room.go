package room

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var _ RedisModel = (*Room)(nil)

type Room struct {
	Id          int64  `redis:"id"`
	Name        string `redis:"name"`
	Desc        string `redis:"desc"`
	OwnerOpenId string `redis:"owner_open_id"`
	connMap     map[int64]*User
	lock        sync.RWMutex
}

func NewRoom(id int64, name, desc string) *Room {
	return &Room{
		Id:      id,
		Name:    name,
		Desc:    desc,
		connMap: make(map[int64]*User),
	}
}

func (r *Room) Fields() map[string]interface{} {
	return map[string]interface{}{
		"id":            r.Id,
		"name":          r.Name,
		"desc":          r.Desc,
		"owner_open_id": r.OwnerOpenId,
	}
}

func (r *Room) Key() string {
	return fmt.Sprintf(rkeyRoomInfo, r.Id)
}

func (r *Room) Save(rc redis.Cmdable) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rc.HMSet(ctx, r.Key(), r.Fields()).Err(); err != nil {
		return err
	}
	return rc.Expire(ctx, r.Key(), 1*time.Hour).Err()
}
