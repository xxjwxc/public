package weixin

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/usthooz/gutil"
	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mycache"
	"github.com/xxjwxc/public/myglobal"
	"github.com/xxjwxc/public/myhttp"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"
)

const (
	_getTicket       = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=wx_card&access_token="
	_getJsurl        = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token="
	_getToken        = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="
	_getSubscribe    = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token="
	_getTempMsg      = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token="
	_createMenu      = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token="
	_deleteMenu      = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token="
	_sendCustom      = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="
	_sendFreepublish = "https://api.weixin.qq.com/cgi-bin/freepublish/batchget?access_token="
	_setGuideConfig  = "https://api.weixin.qq.com/cgi-bin/guide/setguideconfig?access_token="
	_setGetMaterial  = "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token="
	_getUser         = "https://api.weixin.qq.com/cgi-bin/user/get?access_token="
	_cacheToken      = "wx_access_token"
	_cacheTicket     = "weixin_card_ticket"
)

// GetAccessToken 获取微信accesstoken
// 获取登录凭证
func (_wx *wxTools) GetAccessToken() (accessToken string, err error) {
	//先从缓存中获取 access_token
	cache := mycache.NewCache(_cacheToken)
	err = cache.Value(_cacheToken, &accessToken)
	if err == nil {
		return
	}

	var url = _getToken + _wx.wxInfo.AppID + "&secret=" + _wx.wxInfo.AppSecret

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
		if len(accessToken) == 0 {
			mylog.Error(js)
			return
		}
		//保存缓存
		cache.Add(_cacheToken, &accessToken, time.Duration(7000)*time.Second)
		//------------------end
	}
	//----------------------end

	return
}

// clearAccessTokenCache 清除accesstoken缓存
func (_wx *wxTools) clearAccessTokenCache() error {
	//先从缓存中获取 access_token
	cache := mycache.NewCache(_cacheToken)
	return cache.Delete(_cacheToken)
}

// GetAPITicket 获取微信卡券ticket
func (_wx *wxTools) GetAPITicket() (ticket string, err error) {
	//先从缓存中获取
	cache := mycache.NewCache(_cacheTicket)
	err = cache.Value(_cacheTicket, &ticket)
	if err == nil {
		return
	}

	accessToken, e := _wx.GetAccessToken()
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

	return
}

// GetJsTicket 获取微信js ticket
func (_wx *wxTools) GetJsTicket() (ticket string, err error) {
	//先从缓存中获取
	cache := mycache.NewCache("weixin_js_ticket")
	err = cache.Value("base", &ticket)
	if err == nil {
		return
	}
	err = nil

	accessToken, e := _wx.GetAccessToken()
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

	return
}

// SendTemplateMsg 发送订阅消息
func (_wx *wxTools) SendTemplateMsg(msg TempMsg) bool {
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		mylog.Error(err)
		return false
	}

	bo, _ := json.Marshal(msg)
	resb, _ := myhttp.OnPostJSON(_getSubscribe+accessToken, string(bo))

	var res ResTempMsg
	json.Unmarshal(resb, &res)
	return res.Errcode == 0
}

// SendWebTemplateMsg 发送订阅消息
func (_wx *wxTools) SendWebTemplateMsg(msg TempWebMsg) error {
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		mylog.Errorf("SendWebTemplateMsg error: openid:%v,err:%v", msg.Touser, err)
		return err
	}

	bo, _ := json.Marshal(msg)
	resb, _ := myhttp.OnPostJSON(_getTempMsg+accessToken, string(bo))

	var res ResTempMsg
	json.Unmarshal(resb, &res)
	b := res.Errcode == 0
	if !b { // try again
		_wx.clearAccessTokenCache()
		accessToken, err = _wx.GetAccessToken()
		if err != nil {
			mylog.Error(err)
			return err
		}
		resb, _ = myhttp.OnPostJSON(_getTempMsg+accessToken, string(bo))
		json.Unmarshal(resb, &res)
		b = res.Errcode == 0
		if !b {
			if res.Errcode == 43004 {
				return message.GetError(message.Unfollow)
			}
			mylog.Errorf("SendWebTemplateMsg error: openid:%v,res:%v", msg.Touser, res)
		}
	}
	return nil
}

// CreateMenu 创建自定义菜单
func (_wx *wxTools) CreateMenu(menu WxMenu) error { // 创建自定义菜单
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		return err
	}
	bo, _ := json.Marshal(menu)
	resb, _ := myhttp.OnPostJSON(_createMenu+accessToken, string(bo))

	var res ResTempMsg
	json.Unmarshal(resb, &res)
	b := res.Errcode == 0
	if !b {
		return fmt.Errorf("SendWebTemplateMsg error: res:%v", res)
	}

	return nil
}

// DeleteMenu 删除自定义菜单
func (_wx *wxTools) DeleteMenu() error { // 创建自定义菜单
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		return err
	}
	resb, _ := myhttp.OnPostJSON(_deleteMenu+accessToken, string(""))

	var res ResTempMsg
	json.Unmarshal(resb, &res)
	b := res.Errcode == 0
	if !b {
		return fmt.Errorf("SendWebTemplateMsg error: res:%v", res)
	}

	return nil
}

// SendCustomMsg 发送客服消息
func (_wx *wxTools) SendCustomMsg(msg CustomMsg) error {
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		return err
	}
	bo, _ := json.Marshal(msg)
	resb, _ := myhttp.OnPostJSON(_sendCustom+accessToken, string(bo))

	var res ResTempMsg
	json.Unmarshal(resb, &res)
	b := res.Errcode == 0
	if !b {
		return fmt.Errorf("SendWebTemplateMsg error: res:%v", res)
	}

	return nil
}

// SetGuideConfig 快捷回复与关注自动回复
func (_wx *wxTools) SetGuideConfig(guideConfig GuideConfig) error {
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		return err
	}
	bo, _ := json.Marshal(guideConfig)
	resb, _ := myhttp.OnPostJSON(_setGuideConfig+accessToken, string(bo))

	var res ResTempMsg
	json.Unmarshal(resb, &res)
	b := res.Errcode == 0
	if !b {
		return fmt.Errorf("SetGuideConfig error: res:%v", res)
	}

	return nil
}

// GetJsSign js-sdk 授权
func (_wx *wxTools) GetJsSign(url string) (*WxJsSign, error) {
	jsTicket, err := _wx.GetJsTicket()
	if err != nil {
		return nil, err
	}
	// splite url
	urlSlice := strings.Split(url, "#")
	jsSign := &WxJsSign{
		Appid:     _wx.wxInfo.AppID,
		Noncestr:  gutil.RandString(16),
		Timestamp: strconv.FormatInt(time.Now().UTC().Unix(), 10),
		Url:       urlSlice[0],
	}
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", jsTicket, jsSign.Noncestr, jsSign.Timestamp, url)))
	jsSign.Signature = fmt.Sprintf("%x", h.Sum(nil))
	return jsSign, nil
}

// GetAllOpenId  获取用户列表
func (_wx *wxTools) GetAllOpenId() ([]string, error) {
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		return nil, err
	}

	nextOpenid := ""
	var out []string

	for {
		url := _getUser + accessToken
		if len(nextOpenid) > 0 {
			url += fmt.Sprintf("&next_openid=%v", nextOpenid)
		}
		var tmp WxGetUser
		b := myhttp.SendGet(url, "", &tmp)
		if !b {
			return nil, fmt.Errorf("GetAllOpenId error: res:%v", b)
		}
		out = append(out, tmp.Data.Openid...)
		nextOpenid = tmp.Data.NextOpenid
		if len(nextOpenid) == 0 {
			break
		}
	}

	return out, nil
}

// GetFreepublish  获取成功发布列表
func (_wx *wxTools) GetFreepublish(max int64) (out []FreepublishiInfo, err error) {
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		return nil, err
	}
	req := FreepublishiInfoReq{
		Offset:    0,
		Count:     20,
		NoContent: 1,
	}
	if req.Count > max {
		req.Count = max
	}
	for {
		bo, _ := json.Marshal(req)
		resb, _ := myhttp.OnPostJSON(_sendFreepublish+accessToken, string(bo))
		var res FreepublishiInfoResp
		json.Unmarshal(resb, &res)
		if res.ItemCount == 0 {
			break
		}

		for _, v := range res.Item {
			var item FreepublishiInfo
			item.ArticleId = v.ArticleId
			item.UpdateTime = v.UpdateTime
			for _, v := range v.Content.NewsItem {
				item.Title = v.Title
				item.Author = v.Author
				item.Digest = v.Digest
				item.ContentSourceUrl = v.ContentSourceUrl
				item.Url = v.Url
				item.ThumbMediaId = v.ThumbMediaId
				item.IsDeleted = v.IsDeleted
				out = append(out, item)
			}
		}
		if len(out) >= int(max) {
			break
		}
		req.Offset = req.Count
		req.Count += 20
	}

	return out, nil
}

// GetMaterial  获取素材地址
func (_wx *wxTools) GetMaterial(mediaId string) (string, error) {
	accessToken, err := _wx.GetAccessToken()
	if err != nil {
		return "", err
	}
	req := MediaIdReq{
		MediaId: mediaId,
	}

	bo, _ := json.Marshal(req)
	resb, _ := myhttp.OnPostJSON(_setGetMaterial+accessToken, string(bo))
	var res MediaResp
	fmt.Println(string(resb))
	err = json.Unmarshal(resb, &res)
	if err != nil {
		id := fmt.Sprintf("/file/img/%v.jpg", myglobal.GetNode().GetID())
		fileName := path.Join(tools.GetCurrentDirectory(), id)
		tools.WriteFileEx(fileName, resb, true)
		return id, nil
	}
	if len(res.DownUrl) > 0 {
		return res.DownUrl, nil
	}
	if len(res.NewsItemS) > 0 {
		return res.NewsItemS[0].Url, nil
	}

	return "", nil
}
