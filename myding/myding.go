package myding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/xxjwxc/public/myhttp"
	"github.com/xxjwxc/public/tools"
)

type robot struct {
	accessToken string
	secret      string
}
type MsgInfo struct {
	Msgtype  string `json:"msgtype"`
	At       At
	Markdown *Markdown `json:"markdown,omitempty"`
	Text     *Text     `json:"text,omitempty"`
	Link     *Link     `json:"link,omitempty"`
}
type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type Link struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

type Resp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func NewRobot(accessToken, secret string) *robot {
	return &robot{accessToken: accessToken, secret: secret}
}

func (r *robot) SendMsg(msg MsgInfo) error {
	if msg.Text != nil {
		msg.Msgtype = "text"
	} else if msg.Markdown != nil {
		msg.Msgtype = "markdown"
	} else if msg.Link != nil {
		msg.Msgtype = "link"
	}
	timestamp := time.Now().UnixMilli()
	out, err := myhttp.OnPostJSON(fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%v&timestamp=%v&sign=%v", r.accessToken, timestamp, r.Sign(timestamp)), tools.JSONDecode(msg))
	if err != nil {
		return err
	}
	var resp Resp
	tools.JSONEncode(string(out), &resp)
	if resp.Errcode != 0 {
		return fmt.Errorf("ding send err:%v", resp.Errmsg)
	}
	return nil
}

func (r *robot) Sign(timestamp int64) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, r.secret)
	hash := hmac.New(sha256.New, []byte(r.secret))
	hash.Write([]byte(stringToSign))
	signData := hash.Sum(nil)
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signData))
}
