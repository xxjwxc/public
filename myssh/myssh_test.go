package myssh

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	c, err := New("175.24.103.30", "ubuntu", "qwer@1234", 22)
	if err != nil {
		fmt.Println("err", err)
	}

	output, err := c.Run("ls")
	fmt.Printf("%v\n%v", output, err) // 返回字符串

	time.Sleep(1 * time.Second)

	// c.RunTerminal("top") 交互式

	// time.Sleep(1 * time.Second)

	c.Terminal() // 进入
}
