package myredis

import (
	"fmt"
	"testing"
	"time"
)

func Test_cache(t *testing.T) {
	conf := InitRedis(WithAddr("192.155.1.150:6379"), WithClientName(""),
		WithPool(2, 2),
		WithTimeout(10*time.Second), WithReadTimeout(10*time.Second), WithWriteTimeout(10*time.Second),
		WithPwd("Niren1015"), WithGroupName("gggg"), WithDB(0))
	//获取
	res, err := NewRedis(conf)

	fmt.Println(err)
	aaa := "ccccc"
	res.Add("aaaa", aaa, 20*time.Second)
	res.Close()
	res.Add("bbbb", aaa, 0)
	res.Close()
	fmt.Println(res.Ping())

	fmt.Print(res.IsExist("aaaa"))
	fmt.Print(res.GetKeyS("*"))
	fmt.Println(res.Clear())

	fmt.Println(res.Delete("aaaa"))

	var tt string
	res.Value("bbbb", &tt)

	var ww []int32
	res.Add("cccc", []int32{1, 2, 3, 4}, 0)
	res.Value("cccc", &ww)

	// time.Sleep(20 * time.Second)
	// fmt.Print(res.IsExist("aaaa"))

	// fmt.Println(res.Clear())

	// fmt.Println(tt)
	return
}
