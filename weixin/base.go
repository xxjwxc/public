package weixin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"public/mycache"
	"public/mylog"
	"time"

	"github.com/silenceper/wechat"
)

const (
	GETTICKETURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=wx_card&access_token="
	GETJSURL     = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token="
)

//获取微信accesstoken
func GetAccessToken() (access_token string, err error) {
	//先从缓存中获取
	cache := mycache.OnGetCache("weixin_token")
	var tp interface{}
	tp, b := cache.Value("base")
	if b {
		access_token = tp.(string)
	} else {
		wc := wechat.NewWechat(&cfg)
		access_token, err = wc.GetAccessToken()
		//保存缓存
		cache.Add("base", access_token, 7000*time.Second)
	}
	return
}

//获取微信卡券ticket
func GetApiTicket() (ticket string, err error) {
	//先从缓存中获取
	cache := mycache.OnGetCache("weixin_card_ticket")
	var tp interface{}
	tp, b := cache.Value("base")
	if b {
		ticket = tp.(string)
	} else {
		access_token, e := GetAccessToken()
		if e != nil {
			mylog.Error(e)
			err = e
			return
		}
		var url = GETTICKETURL + access_token

		resp, e1 := http.Get(url)
		if e1 != nil {
			mylog.Error(e1)
			err = e1
			return
		}
		defer resp.Body.Close()
		body, e2 := ioutil.ReadAll(resp.Body)
		if e2 != nil {
			mylog.Error(e2)
			err = e2
			return
		}
		var result ApiTicket
		json.Unmarshal(body, &result)
		ticket = result.Ticket
		//保存缓存
		cache.Add("base", ticket, 7000*time.Second)
	}
	return
}

//获取微信js ticket
func GetJsTicket() (ticket string, err error) {
	//先从缓存中获取
	cache := mycache.OnGetCache("weixin_js_ticket")
	var tp interface{}
	tp, b := cache.Value("base")
	if b {
		ticket = tp.(string)
	} else {
		access_token, e := GetAccessToken()
		if e != nil {
			mylog.Error(e)
			err = e
			return
		}
		var url = GETJSURL + access_token

		resp, e1 := http.Get(url)
		if e1 != nil {
			mylog.Error(e1)
			err = e1
			return
		}
		defer resp.Body.Close()
		body, e2 := ioutil.ReadAll(resp.Body)
		if e2 != nil {
			mylog.Error(e2)
			err = e2
			return
		}
		var result ApiTicket
		json.Unmarshal(body, &result)
		ticket = result.Ticket
		//保存缓存
		cache.Add("base", ticket, 7000*time.Second)
	}
	return
}
