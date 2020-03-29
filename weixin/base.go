package weixin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/xxjwxc/public/mycache"
	"github.com/xxjwxc/public/myhttp"
	"github.com/xxjwxc/public/mylog"
)

const (
	_getTicket    = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=wx_card&access_token="
	_getJsurl     = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token="
	_getToken     = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="
	_getSubscribe = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token="
	_cacheToken   = "wx_access_token"
	_cacheTicket  = "weixin_card_ticket"
)

// GetAccessToken 获取微信accesstoken

//获取登录凭证
func GetAccessToken() (accessToken string, err error) {
	//先从缓存中获取 access_token
	cache := mycache.OnGetCache(_cacheToken)
	var tp interface{}
	var b bool
	tp, b = cache.Value(_cacheToken)
	if b {
		accessToken = *(tp.(*string))
	} else {
		var url = _getToken + wxInfo.AppID + "&secret=" + wxInfo.AppSecret

		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		//注入client ip
		js, err := simplejson.NewJson(body)
		if err == nil {
			accessToken, _ = js.Get("access_token").String()
			//保存缓存
			cache.Add(_cacheToken, &accessToken, time.Duration(7000)*time.Second)
			//------------------end
		}
		//----------------------end
	}
	//---------------获取 access_token --------end

	return
}

// GetAPITicket 获取微信卡券ticket
func GetAPITicket() (ticket string, err error) {
	//先从缓存中获取
	cache := mycache.OnGetCache(_cacheTicket)
	var tp interface{}
	tp, b := cache.Value(_cacheTicket)
	if b {
		ticket = tp.(string)
	} else {
		accessToken, e := GetAccessToken()
		if e != nil {
			mylog.Error(e)
			err = e
			return
		}
		var url = _getTicket + accessToken

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
		var result APITicket
		json.Unmarshal(body, &result)
		ticket = result.Ticket
		//保存缓存
		cache.Add(_cacheTicket, ticket, 7000*time.Second)
	}
	return
}

// GetJsTicket 获取微信js ticket
func GetJsTicket() (ticket string, err error) {
	//先从缓存中获取
	cache := mycache.OnGetCache("weixin_js_ticket")
	var tp interface{}
	tp, b := cache.Value("base")
	if b {
		ticket = tp.(string)
	} else {
		accessToken, e := GetAccessToken()
		if e != nil {
			mylog.Error(e)
			err = e
			return
		}
		var url = _getJsurl + accessToken

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
		var result APITicket
		json.Unmarshal(body, &result)
		ticket = result.Ticket
		//保存缓存
		cache.Add("base", ticket, 7000*time.Second)
	}
	return
}

// 发送订阅消息
func SendTemplateMsg(msg TempMsg) bool {
	accessToken, err := GetAccessToken()
	if err != nil {
		mylog.Error(err)
		return false
	}

	bo, _ := json.Marshal(msg)
	resb := myhttp.OnPostJSON(_getSubscribe+accessToken, string(bo))

	var res ResTempMsg
	json.Unmarshal(resb, &res)
	return res.Errcode == 0
}
