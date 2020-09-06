package mykdniao

const (
	// EbusinessOrderHandleUrl 在途监控API
	EbusinessOrderHandleUrl = "http://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
)

// KdnTrace 快递线路图
type KdnTrace struct {
	Action        string `json:"Action"`        // 当前状态
	AcceptStation string `json:"AcceptStation"` // 描述
	AcceptTime    string `json:"AcceptTime"`    // 时间
	Location      string `json:"Location"`      // 当前城市
}

// KdnLogistics 快递鸟返回信息
type KdnLogistics struct {
	LogisticCode string      `json:"LogisticCode"` // 物流运单号
	ShipperCode  string      `json:"ShipperCode"`  // 快递公司编码
	Success      bool        `json:"Success"`      // 成功与否
	Reason       string      `json:"Reason"`       // 失败原因
	State        string      `json:"State"`        // 物流状态：2-在途中,3-签收,4-问题件
	StateEx      string      `json:"StateEx"`      // 增值物流状态： 1-已揽收， 2-在途中， 201-到达派件城市， 202-派件中， 211-已放入快递柜或驿站， 3-已签收， 311-已取出快递柜或驿站， 4-问题件， 401-发货无信息， 402-超时未签收， 403-超时未更新， 404-拒收（退件）， 412-快递柜或驿站超时未取
	Location     string      `json:"Location"`     // 增值所在城市
	Traces       []*KdnTrace `json:"Traces"`       //
}

type kdnLogisticsReq struct {
	LogisticCode string `json:"LogisticCode"`           // 物流单号
	ShipperCode  string `json:"ShipperCode"`            // 快递公司编码
	CustomerName string `json:"CustomerName,omitempty"` // 寄件人或收件人的手机号后四位数字
}
