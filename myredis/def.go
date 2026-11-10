package myredis

import (
	"errors"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const keepTTL = -1

// ErrLockFail 加锁失败(key 已被占用)。
// SET key value NX 未命中时 Redis 返回 nil,redigo 会转成 redis.ErrNil,
// 这里统一包装成明确的业务错误,便于调用方区分"锁被占用"和"连接异常"。
var ErrLockFail = errors.New("lock fail")

// MyRedis redis配置项
type MyRedis struct {
	conf *redisOptions
	// con  redis.Conn
	pool *redis.Pool
	mtx  sync.Mutex
	// once sync.Once
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
