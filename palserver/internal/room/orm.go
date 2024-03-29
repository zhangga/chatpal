package room

type RedisModel interface {
	Key() string
	Fields() map[string]interface{}
}
