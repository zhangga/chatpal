package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	certPath string
	keyPath  string
	address  string
)

func init() {
	flag.StringVar(&certPath, "cert", "", "path to SSL/TLS certificate file")
	flag.StringVar(&keyPath, "key", "", "path to SSL/TLS private key file")
	flag.StringVar(&address, "a", ":10008", "address to use")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	router := gin.Default()
	router.GET("/gm", func(c *gin.Context) {
		c.File("index.html")
	})
	// Handle WebSocket connections
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		log.Println("New connection: ", conn.RemoteAddr().String())
		if err != nil {
			// panic(err)
			log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		for {
			// Read message from client
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				// panic(err)
				log.Printf("%s, error while reading message\n", err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
				break
			}

			// Echo message back to client
			err = conn.WriteMessage(messageType, p)
			if err != nil {
				// panic(err)
				log.Printf("%s, error while writing message\n", err.Error())
				c.AbortWithError(http.StatusInternalServerError, err)
				break
			}
		}
	})

	if certPath == "" || keyPath == "" {
		log.Println("Warning: SSL/TLS certificate and/or private key file not provided. Running server unsecured.")
		err := router.Run(address)
		if err != nil {
			panic(err)
		}
	} else {
		err := router.RunTLS(address, certPath, keyPath)
		if err != nil {
			panic(err)
		}
	}
}
