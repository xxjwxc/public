package mykdniao

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"

	"github.com/xxjwxc/public/mylog"

	"github.com/xxjwxc/public/myhttp"
	"github.com/xxjwxc/public/tools"
)

/*
	快递鸟相关接口
*/

type kdniao struct {
	eBusinessID string
	appKey      string
}

// New 创建一个快递鸟接口
/*
eBusinessID: 商户ID
AppKey: key
*/
func New(eBusinessID, appKey string) *kdniao {
	return &kdniao{eBusinessID: eBusinessID, appKey: appKey}
}

// GetLogisticsTrack 在途监控（获取物流轨迹）
/*
logisticCode:物流单号
shipperCode:快递公司编码
customerName: 当顺丰单号查询时，需要在CustomerName赋值寄件人或收件人的手机号后四位数字
*/
func (k *kdniao) GetLogisticsTrack(logisticCode, shipperCode, customerName string) *KdnLogistics {
	resp := &KdnLogistics{}
	k.post(EbusinessOrderHandleUrl, "8001", &kdnLogisticsReq{LogisticCode: logisticCode, ShipperCode: shipperCode, CustomerName: customerName}, resp)
	if !resp.Success {
		mylog.Errorf("Call Fetch Logistics Failed. Reason: %v", resp.Reason)
		return resp
	}
	return resp
}

func (k *kdniao) post(_url, requestType string, request, resp interface{}) {
	requestData := tools.JSONDecode(request)
	vs := url.Values{}
	vs.Add("RequestData", requestData)
	vs.Add("EBusinessID", k.eBusinessID)
	vs.Add("RequestType", requestType)
	vs.Add("DataSign", k.encrypt(requestData))
	vs.Add("DataType", "2")

	body := myhttp.OnPostForm(_url, vs)

	mylog.Info(string(body))
	if err := json.Unmarshal(body, resp); err != nil {
		mylog.Error(err)
	}
}

func (k *kdniao) encrypt(data string) string {
	data = data + k.appKey
	md5v := md5.Sum(([]byte(data)))
	vv := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(md5v[:])))
	return url.QueryEscape(vv)
}
