package weixin

import (
	"crypto/md5"
	"fmt"
	"log"
	"public/message"
	"public/mylog"
	"public/tools"
	"strconv"
	"strings"
	"time"

	wxpay "gopkg.in/go-with/wxpay.v1"
)

const (
	unifiedOrderUrl = "https://api.mch.weixin.qq.com/pay/unifiedorder"  // 统一下单请求URL
	queryOrderUrl   = "https://api.mch.weixin.qq.com/pay/orderquery"    // 统一查询请求URL
	refundUrl       = "https://api.mch.weixin.qq.com/secapi/pay/refund" //退款请求URL
)

const (
	PAY_SUCCESS    = 1  //支付成功
	PAY_REFUND     = 2  //转入退款
	PAY_CLOSED     = 3  //已关闭
	PAY_NOTPAY     = 4  //未支付
	PAY_REVOKED    = 5  //已撤销
	PAY_USERPAYING = 6  //支付中
	PAY_ERROR      = -1 //支付失败
)

/*
	小程序统一下单接口
	open_id:用户唯一标识
	price : 预支付价钱
	price_body ： 支付描述
	order_id ： 商户订单号
*/
func SmallAppUnifiedorder(open_id string, price int, price_body, order_id, client_ip string) message.MessageBody {
	if !tools.CheckParam(open_id, order_id) || price <= 0 { //参数检测
		return message.GetErrorMsg(message.ParameterInvalid)
	}

	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", client.AppId)
	params.SetString("mch_id", client.MchId)
	params.SetString("body", price_body)
	params.SetInt64("total_fee", int64(price*10))
	params.SetString("spbill_create_ip", client_ip)
	params.SetString("notify_url", notify_url)
	params.SetString("trade_type", "JSAPI")
	params.SetString("openid", open_id)
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("out_trade_no", order_id)               // 商户订单号
	params.SetString("sign", client.Sign(params))            // 签名 c.Sign(params)

	log.Println("paramsparams", params)
	// 发送查询企业付款请求
	ret, err := client.Post(unifiedOrderUrl, params, true)
	if err != nil {
		mylog.Error(err)
		msg := message.GetErrorMsg(message.UnknownError)
		return msg
	}
	//-----------------------end

	//ret["order_id"] = order_tbl.Order_id
	fmt.Println(ret)

	if ret["result_code"] == "SUCCESS" { //再次签名
		dd := make(map[string]string)
		dd["timeStamp"] = strconv.FormatInt(tools.GetUtcTime(time.Now()), 10)
		dd["nonceStr"] = tools.GetRandomString(32)
		dd["package"] = "prepay_id=" + ret["prepay_id"]
		dd["signType"] = "MD5"
		dd["paySign"] = "MD5"
		//appId=wxd678efh567hg6787&nonceStr=5K8264ILTKCH16CQ2502SI8ZNMTM67VS&package=prepay_id=&signType=MD5&timeStamp=1490840662&key=qazwsxedcrfvtgbyhnujmikolp111111
		str := "appId=" + pay_appId + "&nonceStr=" + dd["nonceStr"] + "&package=" + dd["package"] + "&signType=MD5&timeStamp=" + dd["timeStamp"] + "&key=" + apiKey
		by := md5.Sum([]byte(str))
		dd["paySign"] = strings.ToUpper(fmt.Sprintf("%x", by))
		dd["order_id"] = order_id

		msg := message.GetSuccessMsg()
		msg.Data = dd
		return msg
	}

	msg := message.GetErrorMsg(message.InValidOp)
	msg.Data = ret
	return msg
}

/*
	统一查询接口
	open_id:用户唯一标识
	order_id ： 商户订单号
*/
func OnSelectData(open_id, order_id string) (int, message.MessageBody) {
	if !tools.CheckParam(open_id, order_id) { //参数检测
		return 0, message.GetErrorMsg(message.ParameterInvalid)
	}

	code := 0

	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", client.AppId)
	params.SetString("mch_id", client.MchId)
	params.SetString("out_trade_no", order_id)               //商户订单号
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("sign", client.Sign(params))            // 签名 c.Sign(params)

	// 发送查询企业付款请求
	ret := make(wxpay.Params)
	var err error
	ret, err = client.Post(queryOrderUrl, params, true)
	if err != nil { //做再次确认
		time.Sleep(time.Second * 1)
		ret, err = client.Post(queryOrderUrl, params, true)
		if err != nil {
			mylog.Error(err)
			msg := message.GetSuccessMsg()
			return code, msg
		}
	}
	//-----------------------end

	msg := message.GetSuccessMsg(message.NormalMessageId)

	/*
		SUCCESS—支付成功
		REFUND—转入退款
		NOTPAY—未支付
		CLOSED—已关闭
		REVOKED—已撤销（刷卡支付）
		USERPAYING--用户支付中
		PAYERROR--支付失败(其他原因，如银行返回失败)
	*/
	if ret["trade_state"] == "SUCCESS" {
		code = PAY_SUCCESS
	} else if ret["trade_state"] == "REFUND" {
		code = PAY_REFUND
	} else if ret["trade_state"] == "CLOSED" {
		code = PAY_CLOSED
	} else if ret["trade_state"] == "NOTPAY" {
		code = PAY_NOTPAY
	} else if ret["trade_state"] == "REVOKED" {
		code = PAY_REVOKED
	} else if ret["trade_state"] == "USERPAYING" {
		code = PAY_USERPAYING
	} else {
		code = PAY_ERROR
	}

	if ret["trade_state"] == "SUCCESS" { //支付成功
		msg.State = true
	} else {
		msg.State = false
	}

	msg.Data = ret
	return code, msg
}

/*
	申请退款
	open_id:用户唯一标识
	order_id ： 商户订单号
	refund_no：商户退款单号
	total_fee: 订单总金额 分
	refund_fee: 退款总金额 分
*/
func RefundPay(open_id, order_id, refund_no string, total_fee, refund_fee int) (bool, message.MessageBody) {
	if !tools.CheckParam(open_id, order_id) { //参数检测
		return false, message.GetErrorMsg(message.ParameterInvalid)
	}
	code := false
	params := make(wxpay.Params)
	// 退款请求参数
	params.SetString("appid", client.AppId)
	params.SetString("mch_id", client.MchId)
	params.SetString("out_trade_no", order_id)               //商户订单号
	params.SetString("out_refund_no", refund_no)             //商户退款单号
	params.SetInt64("total_fee", int64(total_fee))           // 订单总金额（分）
	params.SetInt64("refund_fee", int64(refund_fee))         // 退款金额（分）
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("sign", client.Sign(params))            // 签名 c.Sign(params)

	// 发送申请退款请求
	ret, err := client.Post(refundUrl, params, true)
	if err != nil {
		mylog.Error(err)
		msg := message.GetErrorMsg(message.UnknownError)
		return code, msg
	}
	//-----------------------end

	msg := message.GetSuccessMsg(message.NormalMessageId)

	if ret["result_code"] == "SUCCESS" { //申请成功
		msg.State = true
		code = true
	} else {
		msg.State = false
	}
	msg.Data = ret
	return code, msg
}
