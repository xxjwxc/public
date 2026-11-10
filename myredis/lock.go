package myredis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xxjwxc/public/mylog"
)

func (mc *redisConPool) TryLock(key interface{}, value interface{}, lifeSpan time.Duration) (err error) {
	var args []interface{}
	args = append(args, mc.getKey(key), mc.encodeValue(value))
	if lifeSpan > 0 {
		if usePrecise(lifeSpan) {
			args = append(args, "PX", formatMs(lifeSpan))
		} else {
			args = append(args, "EX", formatSec(lifeSpan))
		}
	} else if lifeSpan == keepTTL {
		args = append(args, "keepttl")
	}

	args = append(args, "NX")

	con := mc.GetRedisClient()
	defer con.Close()
	repy, err := mc.DO(con, "SET", args...)
	res, err := redis.String(repy, err)
	if err != nil {
		// SET ... NX 未命中时 Redis 返回 nil,redigo 转成 redis.ErrNil,
		// 语义上是"锁已被占用",不应与连接异常混为一谈
		if err == redis.ErrNil {
			return ErrLockFail
		}
		return err
	}

	if mc.conf.isLog {
		mylog.Info(res)
	}
	return nil
}

func (mc *redisConPool) Unlock(key interface{}) (err error) {
	return mc.Delete(key)
}
