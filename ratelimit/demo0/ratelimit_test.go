package demo0

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	ratelimit2 "github.com/juju/ratelimit"
	"github.com/xxjwxc/public/myredis"
	ratelimit1 "go.uber.org/ratelimit"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestTimeRate(t *testing.T) {
	limiter := rate.NewLimiter(rate.Every(2*time.Second), 1)
	for i := 0; i < 10; i++ {
		prev := time.Now()
		now := limiter.Reserve()
		if !now.OK() {
			fmt.Println("no")
		}
		fmt.Println(i, prev)
		time.Sleep(time.Second)
	}
}

func TestRateLimit(t *testing.T) {
	limiter := ratelimit1.New(1, ratelimit1.Per(time.Second*2))
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := limiter.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}

func TestRateLimit2(t *testing.T) {
	bucket := ratelimit2.NewBucket(time.Minute*2, 100)
	fmt.Println(bucket.Available())
	fmt.Println(bucket.TakeAvailable(1))
	fmt.Println(bucket.Available())

	for i := 0; i < 120; i++ {
		before := bucket.Available()
		tokenGet := bucket.TakeAvailable(1)
		if tokenGet != 0 {
			fmt.Println("获取到令牌 index=", i+1, "前后数量-> 前：", before, ", 后: ", bucket.Available(), ", tokenGet=", tokenGet)
		} else {
			fmt.Println("未获取到令牌，拒绝", i+1)
		}
		time.Sleep(1 * time.Second)
	}
}

func TestRedis(t *testing.T) {
	//通过go向redis写入数据和读取数据
	//1. 链接到redis
	conf := myredis.InitRedis(myredis.WithAddr("127.0.0.1:6379"), myredis.WithClientName(""),
		myredis.WithPool(2, 2),
		myredis.WithTimeout(10*time.Second), myredis.WithReadTimeout(10*time.Second), myredis.WithWriteTimeout(10*time.Second),
		myredis.WithPwd(""), myredis.WithGroupName("gggg"), myredis.WithDB(0))
	//获取
	conn, err := myredis.NewRedis(conf)

	//2. 通过go向redis写入数据strinf [key-val]
	_, err = conn.Do("Set", "name", "tomjerry 猫猫")
	if err != nil {
		fmt.Println("set err = ", err)
		return
	}

	//3. 通过go 向redis读取数据string [key-val]
	r, err := redis.String(conn.Do("Get", "name"))
	if err != nil {
		fmt.Println("set err = ", err)
		return
	}
	//因为返回r是interface{}
	//因为name对应的值是string ,因此我们需要转换
	//nameString := r.(string)
	fmt.Println("操作OK", r)
}

func TestNewRateLimit(t *testing.T) {
	bucketName, fillInterval, count := "test1", time.Second*30, 10
	_redis.Delete(bucketName)
	//bucket := NewBucket(fillInterval, 10)
	//
	//_redis.Add(bucketName, &bucket, 0) // 添加redis
	//
	limit := NewRateLimit(bucketName, fillInterval, int64(count))
	for i := 0; i < 130; i++ {
		//var limiter Bucket
		//_redis.Value(bucketName, &limiter)
		//fmt.Println("----", limiter.Available())
		before := limit.Available()
		tokenGet := limit.GetAvailable(bucketName, fillInterval, 1)
		if tokenGet != 0 {
			fmt.Println("获取到令牌 index=", i+1, "前后数量-> 前：", before, ", 后: ", limit.Available(), ", tokenGet=", tokenGet)
		} else {
			fmt.Println("未获取到令牌，拒绝", i+1)
		}
		time.Sleep(1 * time.Second)
	}
}
