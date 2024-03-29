package myredis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xxjwxc/public/mylog"
)

type redisConPool struct {
	base
}

func (mc *redisConPool) Destory() {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()

	if mc.pool != nil {
		mc.pool.Close()
	}
}

func (mc *redisConPool) GetRedisClient() redis.Conn {
	if mc.pool == nil { // 创建连接
		mc.mtx.Lock()
		defer mc.mtx.Unlock()
		mc.pool = &redis.Pool{
			MaxIdle:     mc.conf.maxIdle,
			MaxActive:   mc.conf.maxActive,
			IdleTimeout: mc.conf.timeout,
			Dial: func() (redis.Conn, error) {
				return mc.Dial()
			},
			Wait: (mc.conf.maxActive == 0),
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				if err != nil {
					mylog.Errorf("ping redis error: %s", err)
					return err
				}
				return nil
			},
		}
	}

	return mc.pool.Get()
}

// Ping 判断是否能ping通
func (mc *redisConPool) Ping() bool {
	return mc.ping(mc.GetRedisClient())
}

// Add 添加一个缓存 lifeSpan:缓存时间，0表示永不超时
func (mc *redisConPool) Add(key interface{}, value interface{}, lifeSpan time.Duration) error {
	var args []interface{}
	args = append(args, mc.getKey(key), mc.encodeValue(value))
	if lifeSpan > 0 {
		if usePrecise(lifeSpan) {
			args = append(args, "px", formatMs(lifeSpan))
		} else {
			args = append(args, "ex", formatSec(lifeSpan))
		}
	} else if lifeSpan == keepTTL {
		args = append(args, "keepttl")
	}

	con := mc.GetRedisClient()
	defer con.Close()
	repy, err := mc.DO(con, "SET", args...)
	if mc.conf.isLog {
		mylog.Info(redis.String(repy, err))
	}

	if err != nil {
		mylog.Error(err)
	}
	return err
}

// Value 查找一个cache
func (mc *redisConPool) Value(key interface{}, value interface{}) (err error) {
	con := mc.GetRedisClient()
	defer con.Close()
	repy, err := mc.DO(con, "GET", mc.getKey(key))
	if err != nil {
		mylog.Error(err)
		return err
	}
	return mc.decodeValue(repy, value)
}

// IsExist 判断key是否存在
func (mc *redisConPool) IsExist(key interface{}) bool {
	con := mc.GetRedisClient()
	defer con.Close()
	repy, err := mc.DO(con, "EXISTS", mc.getKey(key))
	if err != nil {
		mylog.Error(err)
		return false
	}
	exist, err := redis.Bool(repy, err) // 转化bool格式
	if err != nil {
		mylog.Error(err)
		return false
	}

	return exist
}

// Delete 删除一个cache
func (mc *redisConPool) Delete(key interface{}) error {
	con := mc.GetRedisClient()
	defer con.Close()
	_, err := mc.DO(con, "del", mc.getKey(key))
	if err != nil {
		mylog.Error(err)
		return err
	}
	return err
}

// Clear 清空表內容
func (mc *redisConPool) Clear() error {
	out, err := mc.GetKeyS("*")
	if err != nil {
		return err
	}

	for _, v := range out {
		err = mc.Delete(v)
		if err != nil {
			return err
		}
	}

	return err
}

// GetKeyS 查询所有key
func (mc *redisConPool) GetKeyS(key interface{}) ([]string, error) {
	con := mc.GetRedisClient()
	defer con.Close()
	var keys []string
	repy, err := mc.DO(con, "keys", mc.getKey(key))
	if err != nil {
		mylog.Error(err)
		return keys, err
	}

	switch t := repy.(type) {
	case []interface{}:
		for _, v := range t {
			out, err := redis.String(v, nil)
			if err != nil {
				mylog.Error(err)
			}
			keys = append(keys, mc.fixKeyGroupName(out))
		}
	default:
		return keys, fmt.Errorf("decodeValue err in type not find:%v", t)
	}

	return keys, err
}

func (mc *redisConPool) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	con := mc.GetRedisClient()
	defer con.Close()
	return mc.DO(con, commandName, args...)
}

// Close 关闭一个连接
func (mc *redisConPool) Close() (err error) {
	return nil
}

// Refresh 更新时间
func (mc *redisConPool) Refresh(key interface{}, lifeSpan time.Duration) error {
	expire := "EXPIRE"
	var args []interface{}
	args = append(args, mc.getKey(key))
	if lifeSpan > 0 {
		if usePrecise(lifeSpan) {
			expire = "PEXPIRE"
			args = append(args, formatMs(lifeSpan))
		} else {
			expire = "EXPIRE"
			args = append(args, formatSec(lifeSpan))
		}
	} else if lifeSpan == keepTTL {
		expire = "EXPIRE"
		args = append(args, 2147483647)
	}

	con := mc.GetRedisClient()
	defer con.Close()
	repy, err := mc.DO(con, expire, args...)
	if mc.conf.isLog {
		mylog.Info(redis.String(repy, err))
	}

	if err != nil {
		mylog.Error(err)
	}
	return err
}
