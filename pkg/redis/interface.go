package redis

import "github.com/gomodule/redigo/redis"

type RedisInterface interface {
	Get(key string) (string, error)
	GetRandomKey() (string, error)
	Set(key string, value interface{}) error
	Del(key string) error
	Pool() *redis.Pool
	Close() error
}
