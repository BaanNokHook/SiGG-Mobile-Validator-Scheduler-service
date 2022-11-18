package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisPoolClient struct {
	redisPool *redis.Pool
}

func NewRedisPoolClient(host string, password string, database int) *RedisPoolClient {
	return &RedisPoolClient{
		redisPool: &redis.Pool{
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			MaxActive: 5,
			MaxIdle:   5,
			Wait:      true,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", host)
				if err != nil {
					return nil, err
				}
				if password != "" {
					if _, err := conn.Do("AUTH", password); err != nil {
						conn.Close()
						return nil, err
					}
				}
				if _, err := conn.Do("SELECT", database); err != nil {
					conn.Close()
					return nil, err
				}
				return conn, nil
			},
		},
	}
}

func (r *RedisPoolClient) Get(key string) (string, error) {
	val, err := redis.String(r.redisPool.Get().Do("GET", key))
	return val, err
}

func (r *RedisPoolClient) GetRandomKey() (string, error) {
	val, err := redis.String(r.redisPool.Get().Do("RANDOMKEY"))
	return val, err
}

func (r *RedisPoolClient) Set(key string, value interface{}) error {
	_, err := redis.Bool(r.redisPool.Get().Do("SET", key, value))
	return err
}

func (r *RedisPoolClient) Del(key string) error {
	_, err := redis.Bool(r.redisPool.Get().Do("DEL", key))
	return err
}

func (r *RedisPoolClient) Pool() *redis.Pool {
	return r.redisPool
}

// Close -.
func (r *RedisPoolClient) Close() error {
	//Close connection here
	return r.redisPool.Close()
}
