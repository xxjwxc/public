package myredis

import (
	"fmt"
	"testing"
	"time"
)

func Test_cache(t *testing.T) {
	//获取
	res, _ := NewRedis([]string{"192.155.1.150:6379"}, "Niren1015", "gggg", 0)

	fmt.Println(res.Clear())

	aaa := "ccccc"
	res.Add("aaaa", aaa, 20*time.Second)
	res.Add("bbbb", aaa, 0)
	fmt.Println(res.Delete("aaaa"))

	fmt.Print(res.IsExist("aaaa"))

	var tt string
	res.Value("aaaa", &tt)

	time.Sleep(20 * time.Second)
	fmt.Print(res.IsExist("aaaa"))

	fmt.Println(res.Clear())

	fmt.Println(tt)
	return
}
