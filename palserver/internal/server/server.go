package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/zhangga/chatpal/palserver/internal/config"
	"github.com/zhangga/chatpal/palserver/internal/room"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	*gin.Engine
	confServer  config.ConfServer
	logger      *zap.SugaredLogger
	redisClient redis.Cmdable
	roomManager *room.Manager
}

func New(ctx context.Context, confPath string) *Server {
	// 加载配置
	conf, err := config.LoadBootstrapConfig(confPath)
	if err != nil {
		panic(err)
	}
	room.AppId, room.AppSecret = conf.Server.AppId, conf.Server.AppSecret

	// 初始化日志
	logger, err := initLogger(&conf.Log)
	if err != nil {
		panic(err)
	}

	// 初始化redis
	redisClient, err := initRedisClient(&conf.Redis)
	if err != nil {
		panic(err)
	}

	roomManager := room.NewManager(ctx, logger, redisClient)

	// 启动web服务
	srv := &Server{
		Engine:      gin.Default(),
		confServer:  conf.Server,
		logger:      logger.Sugar(),
		redisClient: redisClient,
		roomManager: roomManager,
	}
	return srv
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) RunWebsocket() error {
	defer func() {
		s.logger.Sync()
	}()

	if s.Engine == nil {
		s.Engine = gin.Default()
	}

	// home page
	s.Engine.GET("", func(c *gin.Context) {
		c.File("web/index.html")
	})
	// Handle WebSocket connections
	s.Engine.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			s.logger.Errorf("error while Upgrading websocket connection: %s", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		s.logger.Debugf("websocket connection established, remote: %s", conn.RemoteAddr())

		// 加入连接管理
		s.roomManager.AddConn(conn, c)
	})

	if len(s.confServer.Cert) == 0 || len(s.confServer.Key) == 0 {
		s.logger.Infof("SSL/TLS certificate and/or private key file not provided. Running server unsecured.")
		return s.Engine.Run(s.confServer.Addr)
	} else {
		s.logger.Infof("Running server secured with SSL/TLS certificate: %s, private key: %s", s.confServer.Cert, s.confServer.Key)
		return s.Engine.RunTLS(s.confServer.Addr, s.confServer.Cert, s.confServer.Key)
	}
}
