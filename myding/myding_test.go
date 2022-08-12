package myding

import (
	"testing"
)

func Test1(t *testing.T) {
	robot := NewRobot("", "")
	robot.SendMsg(MsgInfo{
		Msgtype: "markdown",
		At:      At{IsAtAll: true},
		Markdown: &Markdown{Title: "招聘信息", Text: `## 姓名: jnn
		### 电话: 123456
		### 公司: 公司
		### 微信: wechat
		#### 详情: 打单发 打单发 打单发 打单发 打单发 打单发 打单发 打单发 打单发 打单发 打单发 打单发 打单发`},
	})
}
