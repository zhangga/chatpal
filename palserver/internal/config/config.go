package config

import (
	"time"
)

type ConfBootstrap struct {
	Server ConfServer `json:"server" yaml:"server"`
	Redis  ConfRedis  `json:"redis" yaml:"redis"`
	Log    ConfLog    `json:"log" yaml:"log"`
}

type ConfServer struct {
	Addr      string `json:"addr" yaml:"addr"`
	Cert      string `json:"cert" yaml:"cert"`
	Key       string `json:"key" yaml:"key"`
	AppId     string `json:"appId" yaml:"appId"`
	AppSecret string `json:"appSecret" yaml:"appSecret"`
}

type ConfRedis struct {
	Cluster      bool          `json:"cluster" yaml:"cluster"`
	AddrList     []string      `json:"addrList" yaml:"addrList"`
	Password     string        `json:"password" yaml:"password"`
	PoolSize     int           `json:"poolSize" yaml:"poolSize"`
	MinIdleConns int           `json:"minIdleConns" yaml:"minIdleConns"`
	MaxRetries   int           `json:"maxRetries" yaml:"maxRetries"`
	DialTimeout  time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	ReadTimeout  time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
}

type ConfLog struct {
	FileName        string `json:"fileName" yaml:"fileName"`
	IsShowInConsole bool   `json:"isShowInConsole" yaml:"isShowInConsole"`
	IsJsonEncoder   bool   `json:"isJsonEncoder" yaml:"isJsonEncoder"`
	Level           string `json:"level" yaml:"level"`
	MaxSize         int    `json:"maxSize" yaml:"maxSize"`
	MaxAge          int    `json:"maxAge" yaml:"maxAge"`
}
