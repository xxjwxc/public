package myredis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/xxjwxc/public/dev"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"
)

type base struct {
	*MyRedis
}

func (mc *base) getCtx() context.Context {
	return context.Background()
}

func (mc *base) getKey(key interface{}) string {
	tmp := ""
	if len(mc.conf.groupName) > 0 {
		tmp = fmt.Sprintf("%v:", mc.conf.groupName)
	}
	switch t := key.(type) {
	case []byte:
		return fmt.Sprintf("%v%v", tmp, string(t))
	case string:
		return fmt.Sprintf("%v%v", tmp, t)
	default:
		return fmt.Sprintf("%v%v", tmp, tools.JSONDecode(key))
	}
}

func (mc *base) encodeValue(value interface{}) interface{} {
	switch t := value.(type) {
	case int32, byte, string, bool, int, uint, int8, int16, int64, uint16, uint32, uint64, float32, float64: // 基础类型
		return t
	default:
		return tools.JSONDecode(value)
		// data, _ := serializing.Encode(value)
		// return data
	}
}

func (mc *base) decodeValue(in, out interface{}) (err error) {
	if in == nil {
		return fmt.Errorf("not fond")
	}

	var reply string
	switch t := in.(type) {
	case []byte:
		reply = string(t)
	default:
		return fmt.Errorf("decodeValue err in type not find:%v", t)
	}

	switch o := out.(type) {
	case *string: // string类型
		*o = reply
		return nil
	case *int32:
		i64, err := strconv.ParseInt(reply, 10, 0)
		*o = int32(i64)
		return err
	case *bool:
		b, err := strconv.ParseBool(string(reply))
		*o = b
		return err
	case *int:
		i64, err := strconv.ParseInt(reply, 10, 0)
		*o = int(i64)
		return err
	case *int8:
		i64, err := strconv.ParseInt(reply, 10, 0)
		*o = int8(i64)
		return err
	case *int16:
		i64, err := strconv.ParseInt(reply, 10, 0)
		*o = int16(i64)
		return err
	case *int64:
		i64, err := strconv.ParseInt(string(reply), 10, 64)
		*o = int64(i64)
		return err
	case *uint:
		i64, err := strconv.ParseUint(reply, 10, 0)
		*o = uint(i64)
		return err
	case *uint8:
		i64, err := strconv.ParseUint(reply, 10, 0)
		*o = uint8(i64)
		return err
	case *uint16:
		i64, err := strconv.ParseUint(reply, 10, 0)
		*o = uint16(i64)
		return err
	case *uint32:
		i64, err := strconv.ParseInt(string(reply), 10, 0)
		*o = uint32(i64)
		return err
	case *uint64:
		i64, err := strconv.ParseUint(reply, 10, 64)
		*o = uint64(i64)
		return err
	case *float32:
		f64, err := strconv.ParseFloat(string(reply), 32)
		*o = float32(f64)
		return err
	case *float64: // 基础类型
		f64, err := strconv.ParseFloat(string(reply), 64)
		*o = float64(f64)
		return err
	default:
		tools.JSONEncode(reply, out) // 复杂类型
		return nil
		//return serializing.Decode(t, out)

	}

	// return fmt.Errorf("decodeValue err not match:%v %v", in, out)
}

func (mc *base) ping(con redis.Conn) bool {
	if con == nil {
		return false
	}

	_, err := con.Do("PING")
	if err != nil {
		mylog.Errorf("ping redis error: %s", err)
		return false
	}
	return true
}

// Dial 获取一个链接
func (mc *base) build() (con redis.Conn, err error) {
	err = fmt.Errorf("not fond ")
	index := mc.conf.addrIdex
	len := len(mc.conf.addrs)
	b := false
	for i := 0; i < len; i++ {
		index = (mc.conf.addrIdex + i) % len
		con, err = redis.Dial("tcp", mc.conf.addrs[index], // redis.DialClientName(mc.conf.clientName),
			redis.DialConnectTimeout(mc.conf.timeout), redis.DialDatabase(mc.conf.db),
			redis.DialPassword(mc.conf.pwd), redis.DialReadTimeout(mc.conf.readTimeout), redis.DialWriteTimeout(mc.conf.writeTimeout),
		)
		if err != nil {
			mylog.Error(err)
		}
		if mc.ping(con) {
			b = true
		}
	}
	if b {
		mc.conf.addrIdex = (index + 1) % len
	}
	return
}

// Dial 获取一个链接
func (mc *base) Dial() (redis.Conn, error) {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()
	return mc.build() // 创建连接
}

func (mc *base) DO(con redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	if dev.IsDev() && mc.conf.isLog {
		cmd := commandName
		for _, v := range args {
			cmd += fmt.Sprintf(" %v", v)
			if len(cmd) > 100 {
				cmd = cmd[:100]
				break
			}
		}
		mylog.Infof("redis req :%v \n", cmd)
	}

	if con != nil {
		reply, err = con.Do(commandName, args...)
		// show log
		if dev.IsDev() && mc.conf.isLog {
			tmp := ""
			switch reply := reply.(type) {
			case []byte:
				tmp = string(reply)
			case string:
				tmp = reply
			case int64:
				tmp = fmt.Sprintf("%v", reply)
			}
			if len(tmp) > 100 {
				tmp = tmp[:100]
			}
			mylog.Infof("redis resp:%v,%v \n", tmp, err)

		}
		return
	}
	return nil, fmt.Errorf("con is nil")
}

func (mc *base) fixKeyGroupName(key string) string {
	tmp := ""
	if len(mc.conf.groupName) > 0 {
		tmp = fmt.Sprintf("%v:", mc.conf.groupName)
	}
	return strings.TrimPrefix(key, tmp)
}

// type redisConOlny struct {
// 	base
// 	con redis.Conn
// }

// // Destory 析构
// func (mc *redisConOlny) Destory() {
// 	mc.mtx.Lock()
// 	defer mc.mtx.Unlock()

// 	if mc.con != nil {
// 		err := mc.con.Close()
// 		if err != nil {
// 			mylog.Error(err)
// 		}
// 		mc.con = nil
// 	}
// }

// // GetRedisClient ...
// func (mc *redisConOlny) GetRedisClient() redis.Conn {
// 	if mc.con == nil {
// 		con, _ := mc.Dial()
// 		mc.con = con
// 	}
// 	return mc.con
// }

// // Ping 判断是否能ping通
// func (mc *redisConOlny) Ping() bool {
// 	return mc.ping(mc.GetRedisClient())
// }

// // 判断是否能ping通
// // Add 添加一个缓存 lifeSpan:缓存时间，0表示永不超时
// func (mc *redisConOlny) Add(key interface{}, value interface{}, lifeSpan time.Duration) error {
// 	var args []interface{}
// 	args = append(args, mc.getKey(key), mc.encodeValue(value))
// 	if lifeSpan > 0 {
// 		if usePrecise(lifeSpan) {
// 			args = append(args, "px", formatMs(lifeSpan))
// 		} else {
// 			args = append(args, "ex", formatSec(lifeSpan))
// 		}
// 	} else if lifeSpan == keepTTL {
// 		args = append(args, "keepttl")
// 	}

// 	_, err := mc.DO(mc.GetRedisClient(), "SET", args...)
// 	if err != nil {
// 		mylog.Error(err)
// 	}
// 	return err
// }

// // Value 查找一个cache
// func (mc *redisConOlny) Value(key interface{}, value interface{}) (err error) {
// 	repy, err := mc.DO(mc.GetRedisClient(), "GET", mc.getKey(key))
// 	if err != nil {
// 		mylog.Error(err)
// 		return err
// 	}
// 	return mc.decodeValue(repy, value)
// }

// // IsExist 判断key是否存在
// func (mc *redisConOlny) IsExist(key interface{}) bool {
// 	repy, err := mc.DO(mc.GetRedisClient(), "EXISTS", mc.getKey(key))
// 	if err != nil {
// 		mylog.Error(err)
// 		return false
// 	}
// 	exist, err := redis.Bool(repy, err) // 转化bool格式
// 	if err != nil {
// 		mylog.Error(err)
// 		return false
// 	}

// 	return exist
// }

// // Delete 删除一个cache
// func (mc *redisConOlny) Delete(key interface{}) error {
// 	_, err := mc.DO(mc.GetRedisClient(), "del", mc.getKey(key))
// 	if err != nil {
// 		mylog.Error(err)
// 		return err
// 	}
// 	return err
// }

// // Clear 清空表內容
// func (mc *redisConOlny) Clear() error {
// 	out, err := mc.GetKeyS("*")
// 	if err != nil {
// 		return err
// 	}

// 	for _, v := range out {
// 		err = mc.Delete(v)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return err
// }

// // GetKeyS 查询所有key
// func (mc *redisConOlny) GetKeyS(key interface{}) ([]string, error) {
// 	var keys []string
// 	repy, err := mc.DO(mc.GetRedisClient(), "keys", mc.getKey(key))
// 	if err != nil {
// 		mylog.Error(err)
// 		return keys, err
// 	}

// 	switch t := repy.(type) {
// 	case []interface{}:
// 		for _, v := range t {
// 			out, err := redis.String(v, nil)
// 			if err != nil {
// 				mylog.Error(err)
// 			}
// 			keys = append(keys, mc.fixKeyGroupName(out))
// 		}
// 	default:
// 		return keys, fmt.Errorf("decodeValue err in type not find:%v", t)
// 	}

// 	return keys, err
// }

// func (mc *redisConOlny) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
// 	con := mc.GetRedisClient()
// 	mc.mtx.Lock()
// 	defer mc.mtx.Unlock()
// 	return mc.DO(con, commandName, args...)
// }

// // Close 关闭一个连接
// func (mc *redisConOlny) Close() (err error) {
// 	mc.mtx.Lock()
// 	defer mc.mtx.Unlock()
// 	if mc.con != nil {
// 		err = mc.con.Close()
// 		mc.con = nil
// 	}

// 	return
// }

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

func formatSec(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}
