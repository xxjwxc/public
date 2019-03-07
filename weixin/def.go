package weixin

type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}
type ApiTicket struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Ticket     string `json:"ticket"`
	Expires_in int    `json:"expires_in"`
}