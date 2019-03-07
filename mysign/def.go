package mysign

import "time"

const (
	_sign_data = "_sign_data"
)

//签名地址
type Sign_client_tbl struct {
	Id                int       `gorm:"primary_key"`
	App_key           string    //key
	App_secret        string    //secret
	Expire_time       time.Time //超时时间
	Strict_sign       int       //是否强制验签:0：用户自定义，1：强制
	Strict_verify     int       //是否强制验证:0：用户自定义，1：强制
	Token_expire_time int       //token过期时间
}

//签名必须带的头标记
type Sing_head struct {
	Appid     string `json:"appid,omitempty"`     //appid
	Signature string `json:"signature,omitempty"` //签名
}
