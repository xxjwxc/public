package weixin

import (
	"github.com/xxjwxc/public/tools"

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

var wxInfo WxInfo

var certFile string // 微信支付商户平台证书路径
var keyFile string
var rootcaFile string

// InitWxinfo 初始化配置信息
func InitWxinfo(info WxInfo) {
	wxInfo = info

	certFile = tools.GetModelPath() + certFileLoc
	keyFile = tools.GetModelPath() + keyFileLoc
	rootcaFile = tools.GetModelPath() + rootcaFileLoc

	//使用memcache保存access_token，也可选择redis或自定义cache
	memCache := NewGocache("_winxin_access")
	//配置微信参数
	cfg = wechat.Config{
		AppID:          wxInfo.APIKey,
		AppSecret:      wxInfo.AppSecret,
		Token:          wxInfo.Token,
		EncodingAESKey: wxInfo.EncodingAESKey,
		Cache:          memCache,
	}
	client = wxpay.NewClient(wxInfo.AppID, wxInfo.MchID, wxInfo.APIKey)
	client.WithCert(certFile, keyFile, rootcaFile)
}
