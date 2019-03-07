package weixin

import (
	"data/config"
	"public/tools"

	"github.com/silenceper/wechat"
	wxpay "gopkg.in/go-with/wxpay.v1"
)

const (
	// 微信支付商户平台证书路径
	certFileLoc   = "/cert/apiclient_cert.pem"
	keyFileLoc    = "/cert/apiclient_key.pem"
	rootcaFileLoc = "/cert/rootca.pem"
)

var cfg wechat.Config
var client *wxpay.Client

var pay_appId string // 微信公众平台应用ID
var mchId string     // 微信支付商户平台商户号
var apiKey string    // 微信支付商户平台API密钥
var secret string
var notify_url string
var token string
var encodingAESKey string

var certFile string // 微信支付商户平台证书路径
var keyFile string
var rootcaFile string

func init() {
	wx_info := config.GetWxInfo()
	//配置微信支付参数
	pay_appId = wx_info.AppID
	mchId = wx_info.MchId
	apiKey = wx_info.Key
	secret = wx_info.AppSecret
	notify_url = wx_info.NotifyUrl
	token = wx_info.Token
	encodingAESKey = wx_info.EncodingAESKey

	certFile = tools.GetModelPath() + certFileLoc
	keyFile = tools.GetModelPath() + keyFileLoc
	rootcaFile = tools.GetModelPath() + rootcaFileLoc

	//使用memcache保存access_token，也可选择redis或自定义cache
	memCache := NewGocache("_winxin_access")
	//配置微信参数
	cfg = wechat.Config{
		AppID:          pay_appId,
		AppSecret:      secret,
		Token:          token,
		EncodingAESKey: encodingAESKey,
		Cache:          memCache,
	}
	client = wxpay.NewClient(pay_appId, mchId, apiKey)
	client.WithCert(certFile, keyFile, rootcaFile)
}
