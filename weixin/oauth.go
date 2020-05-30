package weixin

import (
	"io/ioutil"
	"net/http"

	"github.com/xxjwxc/public/mylog"
)

// SmallAppOauth 小程序授权
func (_wx *wxTools) SmallAppOauth(jscode string) string {
	var url = "https://api.weixin.qq.com/sns/jscode2session?appid=" + _wx.wxInfo.AppID + "&secret=" +
		_wx.wxInfo.AppSecret + "&js_code=" + jscode + "&grant_type=authorization_code&trade_type=JSAPI"

	resp, e := http.Get(url)
	if e != nil {
		mylog.Error(e)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(e)
		return ""
	}
	return string(body)
}
