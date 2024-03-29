package room

import (
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func getRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		PoolSize:     5,
		MinIdleConns: 1,
		MaxRetries:   3,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	return client
}

func TestRoom_Save(t *testing.T) {
	room := &Room{
		Id:   0,
		Name: "test",
	}
	client := getRedisClient()
	if err := room.Save(client); err != nil {
		t.Errorf("Save() error = %v", err)
	}
	var roomRead Room
	if err := client.HGetAll(client.Context(), room.Key()).Scan(&roomRead); err != nil {
		t.Errorf("HGetAll() error = %v", err)
	}
}
