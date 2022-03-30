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
	_, err = redis.String(repy, err)
	if err != nil {
		return err
	}

	if mc.conf.isLog {
		mylog.Info(redis.String(repy, err))
	}
	return err
}

func (mc *redisConPool) Unlock(key interface{}) (err error) {
	return mc.Delete(key)
}
