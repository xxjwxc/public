package demo0

import (
	"fmt"
	"github.com/xxjwxc/public/myredis"
	"time"
)

var _redis myredis.RedisDial

func init() {
	_redis = InitRedis()
}

func InitRedis() myredis.RedisDial {
	conf := myredis.InitRedis(myredis.WithAddr("127.0.0.1:6379"), myredis.WithClientName(""),
		myredis.WithPool(2, 2),
		myredis.WithTimeout(10*time.Second), myredis.WithReadTimeout(10*time.Second), myredis.WithWriteTimeout(10*time.Second),
		myredis.WithPwd(""), myredis.WithGroupName("gggg"), myredis.WithDB(0))
	//获取
	conn, err := myredis.NewRedis(conf)
	if err != nil {
		fmt.Printf("Redis init error: %v\n", err.Error())
	}
	return conn
}

// NewRateLimit 创建一个桶
func NewRateLimit(bucketName string, fillInterval time.Duration, capacity int64) *Bucket {
	var limiter *Bucket
	if !_redis.IsExist(bucketName) { // 不存在 创建一个新的桶
		limiter = NewBucket(fillInterval, capacity) // 新的桶

		err := _redis.Add(bucketName, &limiter, 0) // 添加redis
		if err != nil {
			fmt.Printf("Redis add error: %v\n", err.Error())
		}
	} else { // 存在 直接取
		err := _redis.Value(bucketName, &limiter)
		if err != nil {
			fmt.Printf("Redis get value error: %v\n", err.Error())
		}
	}

	return limiter
}

// GetAvailable 取出桶中count个令牌
func (b *Bucket) GetAvailable(bucketName string, fillInterval time.Duration, count int64) int64 {
	if fillInterval != 0 {
		if time.Now().Sub(b.StartTime) - fillInterval > 0 { // 大于周期 清除redis，重新建桶
			err := _redis.Delete(bucketName)
			if err != nil {
				fmt.Printf("Redis delete error: %v\n", err.Error())
			}
			limit := NewRateLimit(bucketName, fillInterval, count)
			return limit.TakeAvailable(1)
		}
	}

	// 正常取令牌
	return b.TakeAvailable(count)
}
