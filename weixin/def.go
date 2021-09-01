package weixin

import wxpay "gopkg.in/go-with/wxpay.v1"

// UserInfo 用户信息
type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

// APITicket ...
type APITicket struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// WxInfo 微信配置信息
type WxInfo struct {
	AppID     string // 微信公众平台应用ID
	AppSecret string // 微信支付商户平台商户号
	APIKey    string // 微信支付商户平台API密钥
	MchID     string // 商户号
	NotifyURL string // 通知地址
	ShearURL  string // 微信分享url
}

// TempMsg 订阅消息头
type TempMsg struct {
	Touser     string                       `json:"touser"`      //	是	接收者（用户）的 openid
	TemplateID string                       `json:"template_id"` //	是	所需下发的模板消息的id
	Page       string                       `json:"page"`        //	否	点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	Data       map[string]map[string]string `json:"data"`        //是	模板内容，不填则下发空模板
}

// ResTempMsg 模版消息返回值
type ResTempMsg struct {
	Errcode int    `json:"errcode"` //
	Errmsg  string `json:"errmsg"`
}

type wxTools struct {
	client     *wxpay.Client
	wxInfo     WxInfo
	certFile   string // 微信支付商户平台证书路径
	keyFile    string
	rootcaFile string
}

// QrcodeRet ...
type QrcodeRet struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//
type wxPostdata struct {
	Scene string `json:"scene"`
	Page  string `json:"page"`
}

//
type wxQrcodedata struct {
	Path  string `json:"path"`  //路径
	Width int    `json:"width"` //二维码宽度
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`  //	网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn    string `json:"expires_in"`    // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` //	用户刷新access_token
	Openid       string `json:"openid"`        // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope        string `json:"scope"`         // 用户授权的作用域，使用逗号（,）分隔
}

// WxUserinfo 微信用户信息
type WxUserinfo struct {
	Openid     string   `json:"openid"`     // 微信用户唯一标识符,,微信用户唯一标识符
	NickName   string   `json:"nickname"`   // 用户昵称
	Sex        int      `json:"sex"`        // 用户的性别
	City       string   `json:"city"`       // 用户所在城市
	Province   string   `json:"province"`   // 用户所在省份
	Country    string   `json:"country"`    // 用户所在国家
	Headimgurl string   `json:"headimgurl"` // 头像地址
	Privilege  []string `json:"privilege"`  // 户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
}
