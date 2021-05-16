package myredis

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const keepTTL = -1

// MyRedis redis配置项
type MyRedis struct {
	conf *redisOptions
	// con  redis.Conn
	pool *redis.Pool
	mtx  sync.Mutex
	once sync.Once
	dial RedisDial
	err  error
}

var _default = &MyRedis{}

// redisOption redisOption
type redisOptions struct {
	timeout      time.Duration
	groupName    string
	pwd          string
	clientName   string
	addrs        []string
	addrIdex     int
	db           int
	readTimeout  time.Duration
	writeTimeout time.Duration

	// redis pool 相关
	maxIdle   int  // 池中空闲连接的最大数目。
	maxActive int  // 池在给定时间分配的最大连接数。当为零时，池中的连接数没有限制。
	isLog     bool // 是否显示日志
}

// Option 功能选项
type Option interface {
	apply(*redisOptions)
}

type optionFunc func(*redisOptions)

func (f optionFunc) apply(o *redisOptions) {
	f(o)
}
