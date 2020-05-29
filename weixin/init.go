package weixin

import (
	"github.com/xxjwxc/public/tools"

	wxpay "gopkg.in/go-with/wxpay.v1"
)

const (
	// 微信支付商户平台证书路径
	certFileLoc   = "/conf/cert/apiclient_cert.pem"
	keyFileLoc    = "/conf/cert/apiclient_key.pem"
	rootcaFileLoc = "/conf/cert/rootca.pem"
)

var client *wxpay.Client

var wxInfo WxInfo

var certFile string // 微信支付商户平台证书路径
var keyFile string
var rootcaFile string

// InitWxinfo 初始化配置信息
func InitWxinfo(info WxInfo) {
	wxInfo = info

	certFile = tools.GetCurrentDirectory() + certFileLoc
	keyFile = tools.GetCurrentDirectory() + keyFileLoc
	rootcaFile = tools.GetCurrentDirectory() + rootcaFileLoc

	client = wxpay.NewClient(wxInfo.AppID, wxInfo.MchID, wxInfo.APIKey)
	client.WithCert(certFile, keyFile, rootcaFile)
}
