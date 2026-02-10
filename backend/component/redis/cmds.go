package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

func HMSet(key string, values map[string]interface{}, expire int) error {
	conn := GetRedisConn()
	defer conn.Close()
	res, err := redis.String(conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(values)...))
	if err != nil {
		return err
	}
	if res != "OK" {
		return errors.New(res)
	}
	conn.Do("EXPIRE", key, expire)
	return nil
}

func HMGet(key string, fields ...string) (map[string]string, error) {
	conn := GetRedisConn()
	defer conn.Close()
	values, err := redis.Strings(conn.Do("HMGET", redis.Args{}.Add(key).AddFlat(fields)...))
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for index, field := range fields {
		result[field] = values[index]
	}
	return result, nil
}

func HLen(key string) (int, error) {
	conn := GetRedisConn()
	defer conn.Close()
	len, err := redis.Int(conn.Do("HLEN", key))
	if err != nil {
		return 0, err
	}
	return len, nil
}
