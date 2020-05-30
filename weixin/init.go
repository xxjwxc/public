package weixin

import (
	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mylog"

	"github.com/xxjwxc/public/tools"

	wxpay "gopkg.in/go-with/wxpay.v1"
)

const (
	// 微信支付商户平台证书路径
	certFileLoc   = "/conf/cert/apiclient_cert.pem"
	keyFileLoc    = "/conf/cert/apiclient_key.pem"
	rootcaFileLoc = "/conf/cert/rootca.pem"
)

// WxTools 微信操作类型
type WxTools interface {
	GetAccessToken() (accessToken string, err error)                                                        // 获取登录凭证
	GetAPITicket() (ticket string, err error)                                                               // 获取微信卡券ticket
	GetJsTicket() (ticket string, err error)                                                                // 获取微信js ticket
	SendTemplateMsg(msg TempMsg) bool                                                                       // 发送订阅消息
	SmallAppOauth(jscode string) string                                                                     // 小程序授权
	SmallAppUnifiedorder(openID string, price int, priceBody, orderID, clientIP string) message.MessageBody // 小程序统一下单接口
	SelectOrder(openID, orderID string) (int, message.MessageBody)                                          // 统一查询接口
	RefundPay(openID, orderID, refundNO string, totalFee, refundFee int) (bool, message.MessageBody)        // 申请退款
	WxEnterprisePay(openID, tradeNO, desc, ipAddr string, amount int) bool                                  // 企业付款
}

// New 新建及 初始化配置信息
func New(info WxInfo) (WxTools, error) {
	t := &wxTools{
		wxInfo:     info,
		certFile:   tools.GetCurrentDirectory() + certFileLoc,
		keyFile:    tools.GetCurrentDirectory() + keyFileLoc,
		rootcaFile: tools.GetCurrentDirectory() + rootcaFileLoc,
		client:     wxpay.NewClient(info.AppID, info.MchID, info.APIKey),
	}
	err := t.client.WithCert(t.certFile, t.keyFile, t.rootcaFile)
	if err != nil {
		mylog.Error(err)
		return nil, err
	}
	return t, nil
}
