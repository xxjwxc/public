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
			args = append(args, "px", formatMs(lifeSpan))
		} else {
			args = append(args, "ex", formatSec(lifeSpan))
		}
	} else if lifeSpan == keepTTL {
		args = append(args, "keepttl")
	}

	con := mc.GetRedisClient()
	defer con.Close()
	repy, err := mc.DO(con, "SETNX", args...)
	if mc.conf.isLog {
		mylog.Info(redis.String(repy, err))
	}

	if err != nil {
		mylog.Error(err)
	}
	return err
}

func (mc *redisConPool) Unlock(key interface{}) (err error) {
	return mc.Delete(key)
}
