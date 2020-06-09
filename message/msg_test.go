package message

import (
	"fmt"
	"testing"
)

type Test struct {
	State bool   `json:"state"`
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func Test_sing(t *testing.T) {
	var test Test
	test.State = true
	test.Error = ""
	fmt.Println(GetSuccessMsg(NormalMessageID))
	fmt.Println(GetErrorStrMsg("默认的返回值"))
}
