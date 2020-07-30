package redis

import (
	"github.com/garyburd/redigo/redis"
)

func InitRedis(addr string) (redis.Conn, error) {
	Redis, err := redis.Dial("tcp", addr)
	if err != nil {
		return Redis, err
	}
	return Redis, nil
}
