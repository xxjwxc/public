package weixin

import (
	"io/ioutil"
	"net/http"

	"github.com/xxjwxc/public/mylog"
)

/*
	小程序授权
*/
func SmallAppOauth(jscode string) string {
	var url = "https://api.weixin.qq.com/sns/jscode2session?appid=" + wxInfo.AppID + "&secret=" +
		wxInfo.AppSecret + "&js_code=" + jscode + "&grant_type=authorization_code&trade_type=JSAPI"

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
