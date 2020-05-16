package myssh

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	c, err := New("127.0.0.1", "ubuntu", "123456", 22)
	if err != nil {
		fmt.Println("err", err)
	}

	output, err := c.Run("free -h")
	fmt.Printf("%v\n%v", output, err) // 返回字符串

	time.Sleep(1 * time.Second)

	// c.RunTerminal("top") 交互式

	// time.Sleep(1 * time.Second)

	c.Terminal() // 进入
}
