package redis

import (
	"github.com/garyburd/redigo/redis"
)

//RedisIsExist 判断是否存在
func RedisIsExist(Redis redis.Conn, key string) bool {
	Isexist, err := redis.Bool(Redis.Do("EXISTS", key))
	if err != nil {
		return false
	} else {
		return Isexist
	}
}

//RedisSetV 设置value
func RedisSetV(Redis redis.Conn, key string, value interface{}) bool {
	_, err := Redis.Do("SET", key, value)
	if err != nil {
		return false
	}
	return true
}

//RedisSetVWithTime 设置value和过期时间
func RedisSetVWithTime(Redis redis.Conn, key string, value interface{}, seconds string) bool {
	_, err := Redis.Do("SET", key, value, "EX", seconds)
	if err != nil {
		return false
	}
	return true
}

//RedisSetVNX 如果不存在则设置value
func RedisSetVNX(Redis redis.Conn, key string, value interface{}) bool {
	_, err := Redis.Do("SETNX", key, value)
	if err != nil {
		return false
	}
	return true
}

//RedisSetVWithTimeNX 如果不存在则设置value和过期时间
func RedisSetVWithTimeNX(Redis redis.Conn, key string, value interface{}, seconds string) bool {
	_, err := Redis.Do("SETNX", key, value, "EX", seconds)
	if err != nil {
		return false
	}
	return true
}

//RedisSetKeyTime 设置过期时间
func RedisSetKeyTime(Redis redis.Conn, key, seconds string) bool {
	n, _ := Redis.Do("EXPIRE", key, seconds)
	if n == int64(1) {
		return true
	}
	return false
}

//RedisGetV 获取value
func RedisGetV(Redis redis.Conn, key string) (string, error) {
	value, err := redis.String(Redis.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

//RedisGetV 获取value
func RedisGetVBytes(Redis redis.Conn, key string) ([]byte, error) {
	value, err := redis.Bytes(Redis.Do("GET", key))
	if err != nil {
		return []byte{}, err
	}
	return value, nil
}

//RedisDelV 删除value
func RedisDelV(Redis redis.Conn, key string) bool {
	_, err := redis.String(Redis.Do("DEL", key))
	if err != nil {
		return false
	}
	return true
}

//RedisListPush 往列表里插入数据
func RedisListPush(Redis redis.Conn, key string, value interface{}) bool {
	_, err := Redis.Do("lpush", key, value)
	if err != nil {
		return false
	}
	return true
}

//RedisListGetWithRange 从列表里获取指定范围的数据
func RedisListGetWithRange(Redis redis.Conn, key, range1, range2 string) (interface{}, bool) {
	values, err := redis.Values(Redis.Do("lrange", key, range1, range2))
	if err != nil {
		return values, false
	}
	return values, true
}

//RedisFlush 清除缓存
func RedisFlush(Redis redis.Conn) bool {
	err := Redis.Flush()
	if err != nil {
		return false
	}
	return true
}
