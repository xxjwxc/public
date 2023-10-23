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
	GetAllOpenId() ([]string, error)                                                                          // 获取所有用户id
	GetFreepublish(max int64) ([]FreepublishiInfo, error)                                                     // 获取成功发布列表，最大条数
	GetMaterial(mediaId string) (string, error)                                                               // 获取素材地址
	GetBlacklist(openid string) ([]string, string, error)                                                     // 获取黑名单列表
	Getuserphonenumber(code string) (string, error)                                                           // 手机号获取凭证
	// --------------------h5------------------------------
	GetWebOauth(code string) (*AccessToken, error)                    // 授权
	GetWebUserinfo(openid, snaccessToken string) (*WxUserinfo, error) // 获取用户信息
	SendWebTemplateMsg(msg TempWebMsg) error                          // 发送公众号模板消息
	CreateMenu(menu WxMenu) error                                     // 创建自定义菜单
	DeleteMenu() error                                                // 删除自定义菜单
	SetGuideConfig(guideConfig GuideConfig) error                     // 快捷回复与关注自动回复

	SendCustomMsg(msg CustomMsg) error             // 发送客服消息
	UploadTmpFile(path, tp string) (string, error) //上传临时文件(tp:媒体文件类型，分别有图片（image）、语音（voice）、视频（video）和缩略图（thumb）)
	// ----------------------------------------------------

	GetJsSign(url string) (*WxJsSign, error) // js-sdk 授权
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
