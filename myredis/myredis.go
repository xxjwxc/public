package myredis

import "github.com/gomodule/redigo/redis"

// RedisDial 操作
type RedisDial interface {
}

// DefaultConf ...
func DefaultConf() *MyRedis {
	if _default.conf == nil {
		InitDefaultRedis()
	}
	return _default
}

// InitDefaultRedis 初始化(必须要优先调用一次)
func InitDefaultRedis(ops ...Option) {
	var tmp = &redisOptions{}
	for _, o := range ops {
		o.apply(tmp)
	}
	if len(tmp.addrs) == 0 {
		tmp.addrs = append(tmp.addrs, ":6379")
	}

	_default.mtx.Lock()
	defer _default.mtx.Unlock()
	_default.conf = tmp
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

	cnf.mtx.Lock()
	defer cnf.mtx.Unlock()
	cnf.conf = tmp
	return cnf
}

// NewRedis 初始化一个(InitDefaultRedis(需要优先调用)) groupName:分组名
func NewRedis(con *MyRedis) (dial RedisDial, err error) {
	if con == nil {
		con = DefaultConf()
	}
	con.once.Do(func() { // 创建连接
		ReDialRedis(con)
	})
	
	return nil, nil
}

// ReDialRedis 重新连接redis
func ReDialRedis(con *MyRedis) {
	con.mtx.Lock()
	defer con.mtx.Unlock()
	if con.conf.maxIdle > 0 { // 创建连接池

		return
	}

	con.dial = 
	// 创建单个连接
	redis.Dial("tcp",con.)
}
