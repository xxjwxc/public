package weixin

import (
	"crypto/md5"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"

	wxpay "gopkg.in/go-with/wxpay.v1"
)

const (
	unifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"  // 统一下单请求URL
	queryOrderURL   = "https://api.mch.weixin.qq.com/pay/orderquery"    // 统一查询请求URL
	refundURL       = "https://api.mch.weixin.qq.com/secapi/pay/refund" //退款请求URL
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

// SmallAppUnifiedorder 小程序统一下单接口
/*
	小程序统一下单接口
	open_id:用户唯一标识
	price : 预支付价钱
	price_body ： 支付描述
	order_id ： 商户订单号
*/
func SmallAppUnifiedorder(openID string, price int, priceBody, orderID, clientIP string) message.MessageBody {
	if !tools.CheckParam(openID, orderID) || price <= 0 { //参数检测
		return message.GetErrorMsg(message.ParameterInvalid)
	}

	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", client.AppId)
	params.SetString("mch_id", client.MchId)
	params.SetString("body", priceBody)
	params.SetInt64("total_fee", int64(price*10))
	params.SetString("spbill_create_ip", clientIP)
	params.SetString("notify_url", wxInfo.NotifyURL)
	params.SetString("trade_type", "JSAPI")
	params.SetString("openid", openID)
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("out_trade_no", orderID)                // 商户订单号
	params.SetString("sign", client.Sign(params))            // 签名 c.Sign(params)

	log.Println("paramsparams", params)
	// 发送查询企业付款请求
	ret, err := client.Post(unifiedOrderURL, params, true)
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
		str := "appId=" + wxInfo.AppID + "&nonceStr=" + dd["nonceStr"] + "&package=" + dd["package"] + "&signType=MD5&timeStamp=" + dd["timeStamp"] + "&key=" + wxInfo.APIKey
		by := md5.Sum([]byte(str))
		dd["paySign"] = strings.ToUpper(fmt.Sprintf("%x", by))
		dd["order_id"] = orderID

		msg := message.GetSuccessMsg()
		msg.Data = dd
		return msg
	}

	msg := message.GetErrorMsg(message.InValidOp)
	msg.Data = ret
	return msg
}

// OnSelectData 统一查询接口
/*
	统一查询接口
	open_id:用户唯一标识
	order_id ： 商户订单号
*/
func OnSelectData(openID, orderID string) (int, message.MessageBody) {
	if !tools.CheckParam(openID, orderID) { //参数检测
		return 0, message.GetErrorMsg(message.ParameterInvalid)
	}

	code := 0

	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", client.AppId)
	params.SetString("mch_id", client.MchId)
	params.SetString("out_trade_no", orderID)                //商户订单号
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("sign", client.Sign(params))            // 签名 c.Sign(params)

	// 发送查询企业付款请求
	ret := make(wxpay.Params)
	var err error
	ret, err = client.Post(queryOrderURL, params, true)
	if err != nil { //做再次确认
		time.Sleep(time.Second * 1)
		ret, err = client.Post(queryOrderURL, params, true)
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

// RefundPay 申请退款
/*
	申请退款
	open_id:用户唯一标识
	order_id ： 商户订单号
	refund_no：商户退款单号
	total_fee: 订单总金额 分
	refund_fee: 退款总金额 分
*/
func RefundPay(openID, orderID, refundNO string, totalFee, refundFee int) (bool, message.MessageBody) {
	if !tools.CheckParam(openID, orderID) { //参数检测
		return false, message.GetErrorMsg(message.ParameterInvalid)
	}
	code := false
	params := make(wxpay.Params)
	// 退款请求参数
	params.SetString("appid", client.AppId)
	params.SetString("mch_id", client.MchId)
	params.SetString("out_trade_no", orderID)                //商户订单号
	params.SetString("out_refund_no", refundNO)              //商户退款单号
	params.SetInt64("total_fee", int64(totalFee))            // 订单总金额（分）
	params.SetInt64("refund_fee", int64(refundFee))          // 退款金额（分）
	params.SetString("nonce_str", tools.GetRandomString(32)) // 随机字符串
	params.SetString("sign", client.Sign(params))            // 签名 c.Sign(params)

	// 发送申请退款请求
	ret, err := client.Post(refundURL, params, true)
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
