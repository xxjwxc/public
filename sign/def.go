package sign

import "time"

const (
	_sign_data = "_sign_data"
)

type Sign_client_tbl struct {
	Id            int       `gorm:"primary_key"`
	App_key       string    //key
	App_secret    string    //secret
	Expire_time   time.Time //超时时间
	Strict_verify int       //是否强制验证:0：用户自定义，1：强制
}
