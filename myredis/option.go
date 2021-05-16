package myredis

import (
	"time"
)

// WithTimeout 设置过期时间
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(o *redisOptions) {
		o.timeout = timeout
	})
}

// WithGroupName 分组名
func WithGroupName(groupName string) Option {
	return optionFunc(func(o *redisOptions) {
		o.groupName = groupName
	})
}

// WithPwd 密码
func WithPwd(pwd string) Option {
	return optionFunc(func(o *redisOptions) {
		o.pwd = pwd
	})
}

// WithAddr 密码
func WithAddr(addr ...string) Option {
	return optionFunc(func(o *redisOptions) {
		o.addrs = append(o.addrs, addr...)
	})
}

// WithDB 数据库地址(index)
func WithDB(db int) Option {
	return optionFunc(func(o *redisOptions) {
		o.db = db
	})
}

// WithReadTimeout 设置读过期时间
func WithReadTimeout(timeout time.Duration) Option {
	return optionFunc(func(o *redisOptions) {
		o.readTimeout = timeout
	})
}

// WithWriteTimeout 设置写过期时间
func WithWriteTimeout(timeout time.Duration) Option {
	return optionFunc(func(o *redisOptions) {
		o.writeTimeout = timeout
	})
}

// WithPool 连接池配置
func WithPool(maxIdle, maxActive int) Option {
	return optionFunc(func(o *redisOptions) {
		o.maxIdle = maxIdle
		o.maxActive = maxActive
	})
}

// WithClientName 指定Redis服务器连接使用的客户端名称
func WithClientName(name string) Option {
	return optionFunc(func(o *redisOptions) {
		o.clientName = name
	})
}

func WithLog(isLog bool) Option{
	return optionFunc(func(o *redisOptions) {
		o.isLog = isLog
	})	
}