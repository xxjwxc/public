package myredis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"
)

// MyRedis 分布式缓存
type MyRedis struct {
	redis.Cmdable
	timeout   time.Duration
	groupName string
}

// NewRedis 初始化一个cache cachename 缓存名字
func NewRedis(addrs []string, pwd, groupName string, timeout time.Duration) (mc *MyRedis, err error) {
	redis.SetLogger(&logger{})
	if len(addrs) <= 1 {
		mc = &MyRedis{
			Cmdable: redis.NewClient(&redis.Options{
				Addr:        addrs[0],
				DialTimeout: timeout,
				Password:    pwd, // no password set
				// DB:       0,  // use default DB
			}),
		}
	} else {
		mc = &MyRedis{
			Cmdable: redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:       addrs,
				Password:    pwd, // no password set
				DialTimeout: timeout,
			}),
		}
	}
	mc.timeout = timeout
	mc.groupName = groupName

	err = mc.Ping(mc.getCtx()).Err()
	if err != nil {
		mylog.Error(err)
	}
	return
}

func (mc *MyRedis) getCtx() context.Context {
	return context.Background()
}

func (mc *MyRedis) getKey(key interface{}) string {
	if len(mc.groupName) > 0 {
		return fmt.Sprintf("%v:%v", mc.groupName, tools.JSONDecode(key))
	}
	return tools.JSONDecode(key)
}

// Add 添加一个缓存 lifeSpan:缓存时间，0表示永不超时
func (mc *MyRedis) Add(key interface{}, value interface{}, lifeSpan time.Duration) (err error) {
	set := mc.Set(mc.getCtx(), mc.getKey(key), tools.JSONDecode(value), lifeSpan)
	mylog.Info(set.Val())
	// redis.Expect(set.Err()).NotTo(HaveOccurred())
	// redis.Expect(set.Val()).To(Equal("OK"))
	err = set.Err()
	if err != nil {
		mylog.Error(err)
	}
	return
}

// Value 查找一个cache
func (mc *MyRedis) Value(key interface{}, value interface{}) (err error) {
	set := mc.Get(mc.getCtx(), mc.getKey(key))

	mylog.Info(set.Val())
	// redis.Expect(set.Err()).NotTo(HaveOccurred())
	// redis.Expect(set.Val()).To(Equal("OK"))
	err = set.Err()
	if err != nil {
		mylog.Error(err)
		return
	}

	tools.JSONEncode(set.Val(), value)
	return
}

// IsExist 判断key是否存在
func (mc *MyRedis) IsExist(key interface{}) bool {
	set := mc.Exists(mc.getCtx(), mc.getKey(key))
	mylog.Info(set.Val())
	// redis.Expect(set.Err()).NotTo(HaveOccurred())
	// redis.Expect(set.Val()).To(Equal("OK"))
	err := set.Err()
	if err != nil {
		mylog.Error(err)
	}

	return set.Val() == 1
}

// Delete 删除一个cache
func (mc *MyRedis) Delete(key interface{}) error {
	set := mc.Del(mc.getCtx(), mc.getKey(key))
	mylog.Info(set.Val())
	// redis.Expect(set.Err()).NotTo(HaveOccurred())
	// redis.Expect(set.Val()).To(Equal("OK"))
	err := set.Err()
	if err != nil {
		mylog.Error(err)
	}

	return err
}

// GetRedisClient 获取原始cache2go操作类
func (mc *MyRedis) GetRedisClient() redis.Cmdable {
	return mc
}

// Clear 清空表內容
func (mc *MyRedis) Clear() error {
	key := "*"
	if len(mc.groupName) > 0 {
		key = fmt.Sprintf("%v:*", mc.groupName)
	}

	set := mc.Pipeline().Do(mc.getCtx(), "KEYS", key)
	mylog.Info(set.Val())
	// redis.Expect(set.Err()).NotTo(HaveOccurred())
	// redis.Expect(set.Val()).To(Equal("OK"))
	err := set.Err()
	if err != nil {
		mylog.Error(err)
	}

	return mc.delete(key)
}

// delete 删除一个cache
func (mc *MyRedis) delete(key string) error {
	set := mc.Del(mc.getCtx(), key)
	mylog.Info(set.Val())
	// redis.Expect(set.Err()).NotTo(HaveOccurred())
	// redis.Expect(set.Val()).To(Equal("OK"))
	err := set.Err()
	if err != nil {
		mylog.Error(err)
	}

	return err
}
