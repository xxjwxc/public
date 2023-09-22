package weixin

import (
	"encoding/json"
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

// GetWebOauth 网页授权
func (_wx *wxTools) GetWebOauth(code string) (*AccessToken, error) {
	var url = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + _wx.wxInfo.AppID + "&secret=" +
		_wx.wxInfo.AppSecret + "&code=" + code + "&grant_type=authorization_code"

	resp, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, e
	}

	var res AccessToken
	json.Unmarshal(body, &res)
	return &res, nil
}

// GetWebUserinfo 获取用户信息
func (_wx *wxTools) GetWebUserinfo(openid string) (*WxUserinfo, error) {
	accessToken, e := _wx.GetAccessToken()
	if e != nil {
		return nil, e
	}

	var url = "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openid + "&lang=zh_CN"

	resp, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, e
	}

	var res WxUserinfo
	json.Unmarshal(body, &res)
	return &res, nil
}

// Getuserphonenumber 获取用户信息
func (_wx *wxTools) Getuserphonenumber(code string) (string, error) { // 手机号获取凭证
	accessToken, e := _wx.GetAccessToken()
	if e != nil {
		return "", e
	}

	var url = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + accessToken

	code := map[string]interface{}{
		"code": req.Code,
	}
	params, _ := json.Marshal(code)
	resp, e := http.Post(url, "Content-Type", bytes.NewBuffer(params))
	if e != nil {
		return "", e
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", e
	}

	var res WxPhoneResp
	json.Unmarshal(body, &res)
	if res.Errcode != 0 {
		return "", fmt.Errorf(res.Errmsg)
	}
	return res.WxPhoneinfo.PhoneNumber, nil
}
