package weixin

import (
	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mylog"

	"github.com/xxjwxc/public/tools"

	wxpay "gopkg.in/go-with/wxpay.v1"
)

// 微信支付商户平台证书路径

// CertFileLoc  cert.pem
var CertFileLoc = "/conf/cert/apiclient_cert.pem"

// KeyFileLoc key.pem
var KeyFileLoc = "/conf/cert/apiclient_key.pem"

// RootcaFileLoc rootca.pem
var RootcaFileLoc = "/conf/cert/rootca.pem"

// WxTools 微信操作类型
type WxTools interface {
	GetAccessToken() (accessToken string, err error)                                                          // 获取登录凭证
	GetAPITicket() (ticket string, err error)                                                                 // 获取微信卡券ticket
	GetJsTicket() (ticket string, err error)                                                                  // 获取微信js ticket
	SendTemplateMsg(msg TempMsg) bool                                                                         // 发送订阅消息
	SmallAppOauth(jscode string) string                                                                       // 小程序授权
	SmallAppUnifiedorder(openID string, price int64, priceBody, orderID, clientIP string) message.MessageBody // 小程序统一下单接口
	SelectOrder(openID, orderID string) (int, message.MessageBody)                                            // 统一查询接口
	RefundPay(openID, orderID, refundNO string, totalFee, refundFee int64) (bool, message.MessageBody)        // 申请退款
	WxEnterprisePay(openID, tradeNO, desc, ipAddr string, amount int) bool                                    // 企业付款
	GetShareQrcode(path string, scene, page string) (ret QrcodeRet)                                           // 获取小程序码
	GetWxQrcode(path, page string, width int) (ret QrcodeRet)                                                 // 获取小程序二维码 （有限个）

	// --------------------h5------------------------------
	GetWebOauth(code string) (*AccessToken, error)                  // 授权
	GetWebUserinfo(openid, accessToken string) (*WxUserinfo, error) // 获取用户信息
	// ----------------------------------------------------
}

// New 新建及 初始化配置信息
func New(info WxInfo) (WxTools, error) {
	t := &wxTools{
		wxInfo:     info,
		certFile:   tools.GetCurrentDirectory() + CertFileLoc,
		keyFile:    tools.GetCurrentDirectory() + KeyFileLoc,
		rootcaFile: tools.GetCurrentDirectory() + RootcaFileLoc,
		client:     wxpay.NewClient(info.AppID, info.MchID, info.APIKey),
	}
	err := t.client.WithCert(t.certFile, t.keyFile, t.rootcaFile)
	if err != nil {
		mylog.Error(err)
		return nil, err
	}
	return t, nil
}
