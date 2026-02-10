package redis

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var once sync.Once
var pool *redis.Pool

func init() {
	once.Do(func() {
		pool = &redis.Pool{
			MaxIdle:     60,                              // 最大连接数量
			MaxActive:   60,                              // 最大活跃连接数，不确定用0
			IdleTimeout: time.Duration(60) * time.Second, // 空闲连接超时时间|查询:CONFIG GET timeout|设置CONFIG SET timeout 65
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", "121.37.146.209:6380")
			},
		}
	})
}

// 连接池中获取连接
func GetRedisConn() redis.Conn {
	return pool.Get()
}

// 将连接放回连接池
func CloseRedisConn(conn redis.Conn) {
	conn.Close()
}

// 分布式单用户锁
func GetLock(key string, time uint) (redis.Conn, bool, error) {
	conn := pool.Get()
	if time == 0 {
		time = 3
	}
	if time > 10 {
		time = 10
	}
	res, err := redis.String(conn.Do("SET", key, "1", "EX", time, "NX"))
	if err != nil {
		return conn, false, err
	}

	return conn, res == "OK", nil
}

// 释放单用户锁
func Unlock(conn redis.Conn, key string) error {
	defer conn.Close()
	iLoop := 3
	for {
		_, err := conn.Do("DEL", key)
		iLoop--
		if err == nil || iLoop <= 0 {
			break
		}
	}
	return nil
}

// 接口限流
func AccessLimit(key string, cycle, limit int) bool {
	if key == "" {
		return false
	}

	conn := pool.Get()
	defer conn.Close()

	used, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return true // Redis 异常时直接通过
	}

	if used > limit {
		ttl, err := redis.Int(conn.Do("TTL", key))
		if err != nil || ttl == -1 {
			conn.Do("EXPIRE", key, cycle)
		}
		return false
	}

	if used == 1 {
		conn.Do("EXPIRE", key, cycle)
	}

	return true
}
