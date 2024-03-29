package server

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/zhangga/chatpal/palserver/internal/config"
	"time"
)

func initRedisClient(conf *config.ConfRedis) (redis.Cmdable, error) {
	if conf == nil || len(conf.AddrList) == 0 {
		return nil, errors.New("redis config is nil or addr list is empty")
	}

	var client redis.Cmdable
	// redis cluster
	if conf.Cluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        conf.AddrList,
			Password:     conf.Password,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConns,
			MaxRetries:   conf.MaxRetries,
			DialTimeout:  conf.DialTimeout,
			ReadTimeout:  conf.ReadTimeout,
			WriteTimeout: conf.WriteTimeout,
		})
		if err := checkRedisCluster(client); err != nil {
			return nil, err
		}
	} else { // redis standalone
		client = redis.NewClient(&redis.Options{
			Addr:         conf.AddrList[0],
			Password:     conf.Password,
			PoolSize:     conf.PoolSize,
			MinIdleConns: conf.MinIdleConns,
			MaxRetries:   conf.MaxRetries,
			DialTimeout:  conf.DialTimeout,
			ReadTimeout:  conf.ReadTimeout,
			WriteTimeout: conf.WriteTimeout,
		})
	}
	// check redis
	if err := checkRedis(client); err != nil {
		return nil, err
	}
	return client, nil
}

func checkRedis(redisClient redis.Cmdable) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// check redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return err
	}
	// check redis script
	cmd := redis.NewScript("redis.replicate_commands(); return 1").Run(ctx, redisClient, nil)
	if _, err := cmd.Result(); err != nil {
		return err
	}
	return nil
}

func checkRedisCluster(redisClient redis.Cmdable) error {
	// check redis cluster connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := redisClient.ClusterSlots(ctx).Result(); err != nil {
		return err
	}
	return nil
}
