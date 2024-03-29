package room

import (
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/zhangga/chatpal/palserver/internal/msg"
	"go.uber.org/zap"
	"sync"
	"time"
)

const messageType = 2

// conn 封装的websocket连接
type conn struct {
	*websocket.Conn
	id          int64
	ctx         *gin.Context
	logger      *zap.SugaredLogger
	redisClient redis.Cmdable
	user        *User
}

// handleMsg 处理消息
func (c *conn) handleMsg(msgType int, buf []byte) error {
	msgId := binary.LittleEndian.Uint16(buf[0:])
	seq := binary.LittleEndian.Uint32(buf[2:])
	handler, ok := msgHandlers[msg.Id(msgId)]
	if !ok {
		return fmt.Errorf("unknown message Id: %d", msgId)
	}
	content := buf[6:]
	resp, err := handler(c, msgType, content)
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	return c.sendMsg(msg.Id(msgId), resp, seq)
}

func (c *conn) sendMsg(id msg.Id, msg msg.Message, seq uint32) error {
	content, err := json.Marshal(msg)
	if err != nil {
		c.logger.Errorf("error while marshaling message, Id: %d, err: %v", id, err)
		return err
	}
	buf := make([]byte, 2+4+len(content))
	binary.LittleEndian.PutUint16(buf[0:], uint16(id))
	binary.LittleEndian.PutUint32(buf[2:], seq)
	copy(buf[6:], content)
	err = c.WriteMessage(messageType, buf)
	if err != nil {
		c.logger.Errorf("error while sending message, Id: %d, err: %v", id, err)
	}
	return err
}

// connManager websocket连接管理
type connManager struct {
	logger *zap.SugaredLogger
	id     int64
	lock   sync.RWMutex
	conns  map[int64]*conn
}

func newConnManager(logger *zap.Logger) connManager {
	return connManager{
		logger: logger.Sugar(),
		id:     time.Now().Unix(),
		conns:  make(map[int64]*conn),
	}
}

func (c *connManager) addConn(wc *websocket.Conn, ctx *gin.Context, rc redis.Cmdable) *conn {
	conn := &conn{
		ctx:         ctx,
		Conn:        wc,
		redisClient: rc,
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.id++
	id := c.id
	conn.id = id
	conn.logger = c.logger.With("connId", id)
	c.conns[id] = conn
	return conn
}

func (c *connManager) removeConn(id int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.conns, id)
}

func (c *connManager) readloopConn(conn *conn) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				c.logger.Errorf("readloopConn panic: %v", r)
			}
		}()
		defer func() {
			c.logger.Debugf("conn closed, remote: %s", conn.RemoteAddr())
			c.removeConn(conn.id)
			conn.Close()
		}()

		for {
			// Read message from client
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.logger.Infof("closed while reading message, conn: %s, err: %v", conn.RemoteAddr(), err)
				} else {
					c.logger.Errorf("error while reading message, conn: %s, err: %v", conn.RemoteAddr(), err)
				}
				//conn.ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			c.logger.Debugf("received message, remote: %s, type: %d, content: %s", conn.RemoteAddr(), messageType, p)
			if err = conn.handleMsg(messageType, p); err != nil {
				c.logger.Errorf("error while handling message, conn: %s, type: %d, err: %v", conn.RemoteAddr(), messageType, err)
			}
		}
	}()
}
