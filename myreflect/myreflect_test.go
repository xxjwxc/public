package myreflect

import (
	"fmt"
	"testing"
)

type ReqTest1 struct {
	AccessToken string `json:"access_token"`                 // access_token
	UserName    string `json:"user_name" binding:"required"` // user name
	Password    string `json:"password"`                     // password
}

func TestJson(t *testing.T) {
	var ts ReqTest1
	ts.UserName = `1111`
	fmt.Println(FindTag(ts, "UserName", "json"))
}
