package main

import (
	"context"
	"flag"
	"github.com/zhangga/chatpal/palserver/internal/server"
)

var confPath string

func init() {
	flag.StringVar(&confPath, "c", "configs/server.yaml", "config file path, eg: -c=configs/server.yaml")
}

func main() {
	srv := server.New(context.Background(), confPath)
	srv.RunWebsocket()
}
