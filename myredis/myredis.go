package myredis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xxjwxc/public/mylog"
)

// RedisDial 操作
type RedisDial interface {
	Destory()                                                                       // 析构
	GetRedisClient() redis.Conn                                                     // 获取一个原始的redis连接
	Ping() bool                                                                     // 判断是否能ping通
	Add(key interface{}, value interface{}, lifeSpan time.Duration) error           // 添加一个元素
	Value(key interface{}, value interface{}) error                                 // 获取一个value
	IsExist(key interface{}) bool                                                   // 判断是否存在
	Delete(key interface{}) error                                                   // 删除一个
	Clear() error                                                                   // 清空
	GetKeyS(key interface{}) ([]string, error)                                      // 查询所有key
	Close() (err error)                                                             // 关闭连接
	Do(commandName string, args ...interface{}) (reply interface{}, err error)      // 一次操作
	TryLock(key interface{}, value interface{}, lifeSpan time.Duration) (err error) //  试着加锁
	Unlock(key interface{}) (err error)                                             // 解锁
}

// DefaultConf ...
func DefaultConf() *MyRedis {
	if _default.conf == nil {
		InitDefaultRedis()
	}
	return _default
}

// InitDefaultRedis 初始化(必须要优先调用一次)
func InitDefaultRedis(ops ...Option) *MyRedis {
	var tmp = &redisOptions{isLog: true}
	for _, o := range ops {
		o.apply(tmp)
	}
	if len(tmp.addrs) == 0 {
		tmp.addrs = append(tmp.addrs, ":6379")
	}

	_default.mtx.Lock()
	defer _default.mtx.Unlock()
	_default.conf = tmp
	return _default
}

// InitRedis 初始化(必须要优先调用一次)
func InitRedis(ops ...Option) *MyRedis {
	var cnf = &MyRedis{}
	var tmp = &redisOptions{}
	for _, o := range ops {
		o.apply(tmp)
	}
	if len(tmp.addrs) == 0 {
		tmp.addrs = append(tmp.addrs, ":6379")
	}

	// cnf.mtx.Lock()
	// defer cnf.mtx.Unlock()
	cnf.conf = tmp
	return cnf
}

// NewRedis 初始化一个(InitDefaultRedis(需要优先调用)) groupName:分组名
func NewRedis(con *MyRedis) (dial RedisDial, err error) {
	if con == nil {
		con = DefaultConf()
	}

	if con.dial != nil { // 清理，关闭连接
		con.dial.Destory()
	}

	// con.mtx.Lock()
	// defer con.mtx.Unlock()

	if con.conf.maxIdle == 0 { // 创建连接池
		con.conf.maxIdle = 1
	}

	// 创建连接池
	con.dial = &redisConPool{
		base: base{MyRedis: con},
	}
	con.pool = &redis.Pool{
		MaxIdle:     con.conf.maxIdle,
		MaxActive:   con.conf.maxActive,
		IdleTimeout: con.conf.timeout,
		Dial: func() (redis.Conn, error) {
			con.mtx.Lock()
			defer con.mtx.Unlock()
			index := con.conf.addrIdex

			len := len(con.conf.addrs)
			con.conf.addrIdex = (index + 1) % len

			return redis.Dial("tcp", con.conf.addrs[index], // redis.DialClientName(con.conf.clientName),
				redis.DialConnectTimeout(con.conf.timeout), redis.DialDatabase(con.conf.db),
				redis.DialPassword(con.conf.pwd), redis.DialReadTimeout(con.conf.readTimeout), redis.DialWriteTimeout(con.conf.writeTimeout),
			)
		},
		Wait: true,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				mylog.Errorf("ping redis error: %s", err)
				return err
			}
			return nil
		},
	}

	return con.dial, nil
}
