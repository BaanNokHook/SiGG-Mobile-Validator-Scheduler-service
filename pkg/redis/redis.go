package redis

import (
	"github.com/gomodule/redigo/redis"
)

type RedisClient struct {
	redisConn redis.Conn
}

func NewRedisClient(host string, password string, database int) *RedisClient {
	conn, err := redis.Dial("tcp", host)
	if password != "" {
		if _, err := conn.Do("AUTH", password); err != nil {
			conn.Close()
		}
	}
	if _, err := conn.Do("SELECT", database); err != nil {
		conn.Close()
	}
	if err != nil {
		panic(err)
	}
	return &RedisClient{
		redisConn: conn,
	}
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := redis.String(r.redisConn.Do("GET", key))
	return val, err
}

func (r *RedisClient) GetRandomKey() (string, error) {
	val, err := redis.String(r.redisConn.Do("RANDOMKEY"))
	return val, err
}

func (r *RedisClient) Set(key string, value interface{}) error {
	_, err := redis.Bool(r.redisConn.Do("SET", key, value))
	return err
}

func (r *RedisClient) Del(key string) error {
	_, err := redis.Bool(r.redisConn.Do("DEL", key))
	return err
}

func (r *RedisClient) Pool() *redis.Pool {
	return nil
}

// Close -.
func (r *RedisClient) Close() error {
	//Close connection here
	return r.redisConn.Close()
}
