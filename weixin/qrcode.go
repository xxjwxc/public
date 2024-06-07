package weixin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/xxjwxc/public/myhttp"
	"github.com/xxjwxc/public/tools"
)

const (
	GETSHAREURL    = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="          //小程序码
	GETQRCODEURL   = "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=" //小程序二维码
	GETH5QRCODEURL = "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="          //h5二维码
)

// GetShareQrcode 获取小程序码
// path：图片保存路径
// scene: 附带参数
// page:小程序页面头部
func (_wx *wxTools) GetShareQrcode(path string, scene, page string) (ret QrcodeRet) {
	accessToken, _ := _wx.GetAccessToken() // 获取access_token
	data := wxPostdata{Scene: scene, Page: page}
	bo, _ := json.Marshal(data)
	resb, _ := myhttp.OnPostJSON(GETSHAREURL+accessToken, string(bo))

	tools.JSONEncode(string(resb), &ret) //错误码45029 最大限制
	if ret.Errcode == 0 {
		ioutil.WriteFile(path, resb, 0666) //写入文件(字节数组)
	}
	return
}

// GetWxQrcode 获取小程序二维码 （有限个）
// path：图片保存路径
// page: 小程序页面pages/index?query=1
// width: 二维码宽度

func (_wx *wxTools) GetWxQrcode(path, page string, width int) (ret QrcodeRet) {
	fmt.Println(path)
	//获取access_token
	accessToken, _ := _wx.GetAccessToken()

	data := wxQrcodedata{Path: page, Width: width}
	bo, _ := json.Marshal(data)
	resb, _ := myhttp.OnPostJSON(GETQRCODEURL+accessToken, string(bo))

	tools.JSONEncode(string(resb), &ret) //错误码45029 最大限制
	if ret.Errcode == 0 {
		ioutil.WriteFile(path, resb, 0666) //写入文件(字节数组)
	}
	return
}

// GetH5Qrcode 生成带参数的二维码
// expireSeconds : 该二维码有效时间，以秒为单位。 最大不超过2592000（即30天），此字段如果不填，则默认有效期为60秒。
// actionName: 二维码类型，QR_SCENE为临时的整型参数值，QR_STR_SCENE为临时的字符串参数值，QR_LIMIT_SCENE为永久的整型参数值，QR_LIMIT_STR_SCENE为永久的字符串参数值
// aceneId: 场景值ID，临时二维码时为32位非0整型，永久二维码时最大值为100000（目前参数只支持1--100000

func (_wx *wxTools) GetH5Qrcode(expireSeconds int, did string, aceneId int) (string, error) { // 手机号获取凭证
	accessToken, e := _wx.GetAccessToken()
	if e != nil {
		return "", e
	}

	req := H5QrcodeReq{
		ExpireSeconds: expireSeconds,
		ActionName:    "QR_STR_SCENE",
		SceneId:       aceneId,
		ActionInfo: ActionInfo{
			Scene: ActionScene{
				SceneStr: did,
			},
		},
	}
	resb, _ := myhttp.OnPostJSON(GETH5QRCODEURL+accessToken, tools.JSONDecode(req))

	resp := H5QrcodeResp{}
	tools.JSONEncode(string(resb), &resp) //错误码45029 最大限制
	if resp.Ticket == "" {
		_wx.clearAccessTokenCache()
		return "", fmt.Errorf("ticket is null")
	}
	return fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", url.QueryEscape(resp.Ticket)), nil
}
