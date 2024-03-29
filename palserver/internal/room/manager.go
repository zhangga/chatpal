package room

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Manager struct {
	ctx         context.Context
	logger      *zap.SugaredLogger
	redisClient redis.Cmdable
	conns       connManager
}

func NewManager(ctx context.Context, logger *zap.Logger, rc redis.Cmdable) *Manager {
	return &Manager{
		ctx:         ctx,
		logger:      logger.Sugar(),
		redisClient: rc,
		conns:       newConnManager(logger),
	}
}

// AddConn 添加websocket连接
func (m *Manager) AddConn(wc *websocket.Conn, ctx *gin.Context) {
	conn := m.conns.addConn(wc, ctx, m.redisClient)
	m.conns.readloopConn(conn)
}
