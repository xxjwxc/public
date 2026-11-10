package weixin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/xxjwxc/public/mylog"
)

// SmallAppOauth 小程序授权
func (_wx *wxTools) SmallAppOauth(jscode string) string {
	var url = "https://api.weixin.qq.com/sns/jscode2session?appid=" + _wx.wxInfo.AppID + "&secret=" +
		_wx.wxInfo.AppSecret + "&js_code=" + jscode + "&grant_type=authorization_code&trade_type=JSAPI"

	resp, e := httpGet(url)
	if e != nil {
		mylog.Error(e)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(err)
		return ""
	}
	return string(body)
}

// GetWebOauth 网页授权
func (_wx *wxTools) GetWebOauth(code string) (*AccessToken, error) {
	var url = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + _wx.wxInfo.AppID + "&secret=" +
		_wx.wxInfo.AppSecret + "&code=" + code + "&grant_type=authorization_code"

	resp, e := httpGet(url)
	if e != nil {
		return nil, e
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res AccessToken
	err = json.Unmarshal(body, &res)
	if err != nil {

		mylog.Error(err, string(body))
		return nil, err
	}
	return &res, nil
}

// GetWebUserinfo 获取用户信息
func (_wx *wxTools) GetWebUserinfo(openid, snaccessToken string) (*WxUserinfo, error) {
	var url = "https://api.weixin.qq.com/sns/userinfo?access_token=" + snaccessToken + "&openid=" + openid + "&lang=zh_CN"

	resp, e := httpGet(url)
	if e != nil {
		return nil, e
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res WxUserinfo
	if err = json.Unmarshal(body, &res); err != nil {
		mylog.Error(err, string(body))
		return nil, err
	}
	return &res, nil
}

// Getuserphonenumber 获取用户信息
func (_wx *wxTools) Getuserphonenumber(code string) (string, error) { // 手机号获取凭证
	accessToken, e := _wx.GetAccessToken()
	if e != nil {
		return "", e
	}

	var url = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + accessToken

	params := map[string]interface{}{
		"code": code,
	}
	req, _ := json.Marshal(params)
	resp, e := httpPost(url, "Content-Type", bytes.NewBuffer(req))
	if e != nil {
		return "", e
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res WxPhoneResp
	if err = json.Unmarshal(body, &res); err != nil {
		mylog.Error(err, string(body))
		return "", err
	}
	if res.Errcode != 0 {
		return "", errors.New(res.Errmsg)
	}
	return res.WxPhoneinfo.PhoneNumber, nil
}
