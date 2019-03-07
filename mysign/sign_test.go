package mysign

import (
	"fmt"
	"public/message"
	"public/tools"
	"testing"
	"time"
)

func Test_sing(t *testing.T) {
	now := time.Now()
	str := "1" + tools.GetTimeString(now)
	str += "1.0001"
	fmt.Println(str)
	ttt := tools.Md5Encoder(str)
	fmt.Println(ttt)
	fmt.Println(OnCheckSign("wwwthings", ttt, 1, now, 1.0001))
	fmt.Println(message.GetSuccessMsg())
}
