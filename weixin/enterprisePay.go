package weixin

import (
	"log"

	"github.com/xxjwxc/public/tools"

	wxpay "gopkg.in/go-with/wxpay.v1"
)

const (
	enterprisePayURL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers" // 查询企业付款接口请求URL
)

// WxEnterprisePay 企业付款
/*
企业付款
open_id:用户唯一标识
trade_no : 商户订单号
desc ： 操作说明
amount：付款金额 分
*/
func WxEnterprisePay(openID, tradeNO, desc, ipAddr string, amount int) bool {
	c := wxpay.NewClient(wxInfo.AppID, wxInfo.MchID, wxInfo.APIKey)

	// 附着商户证书
	err := c.WithCert(certFile, keyFile, rootcaFile)
	if err != nil {
		log.Fatal(err)
	}

	params := make(wxpay.Params)
	nonceStr := tools.GetRandomString(16)
	// 查询企业付款接口请求参数
	params.SetString("mch_appid", c.AppId)        //商户账号appid
	params.SetString("mchid", c.MchId)            //商户号
	params.SetString("nonce_str", nonceStr)       // 随机字符串
	params.SetString("partner_trade_no", tradeNO) // 商户订单号
	params.SetString("openid", openID)            //用户openid
	params.SetString("check_name", "NO_CHECK")    //校验用户姓名选项
	params.SetInt64("amount", int64(amount))      //企业付款金额，单位为分
	params.SetString("desc", desc)                //企业付款操作说明信息。必填。
	params.SetString("spbill_create_ip", ipAddr)

	params.SetString("sign", c.Sign(params)) // 签名

	// 发送查询企业付款请求
	ret, err := c.Post(enterprisePayURL, params, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(ret)
	returnCode := ret.GetString("return_code")
	resultCode := ret.GetString("result_code")
	if returnCode == "SUCCESS" && resultCode == "SUCCESS" {
		return true
	}

	return false
}
