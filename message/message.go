/*
	消息定义接口
*/
package message

import (
	"encoding/json"
	"log"

	"github.com/bitly/go-simplejson"
	"github.com/xxjwxc/public/dev"
	"github.com/xxjwxc/public/tools"
)

const ( //消息id定义
	NormalMessageId   = 0 ////默认的返回值，为0，自增
	ServerMaintenance = 1 //服务器维护中 请稍后再试,服务器出错会回这个信息，暂时不用
	AccountDisabled   = 2 //帐号被禁用
	AppidOverdue      = 3 //appid过期

	UnknownError  = 101 //未知错误
	TokenFailure  = 102 //token失效
	HTMLSuccess   = 200 //成功
	BlockingAcess = 405 //禁止访问( 禁用请求中所指定的方法)

	NewReport = 2001 //新消息
	NewHeart  = 2002 //心跳

	ParameterInvalid         = 1001 //参数无效
	AppidParameterInvalid    = 1002 //appid参数无效
	EncryptCheckError        = 1003 //密文校验失败
	UserNameDoNotExist       = 1004 //用户名不存在或密码错误
	DuplicateKeyError        = 1005 //键值对重复
	InValidOp                = 1007 //无效操作
	NotFindError             = 1006 //未找到
	InValidAuthorize         = 1008 //授权码错误
	HasusedError             = 1009 //已被使用
	HasActvError             = 1010 //已被激活
	ActvFailure              = 1011 //|激活码被禁止使用
	UserExisted              = 1012 //用户已存在
	VerifyTimeError          = 1013 //验证码请求过于平凡
	MailSendFaild            = 1014 //邮箱发送失败
	SMSSendFaild             = 1015 //手机发送失败
	PhoneParameterError      = 1016 //手机号格式有问题
	VerifyError              = 1017 //验证码错误
	UserNotExisted           = 1018 //用户不存在
	TopicExisted             = 1019 //topic已经存在
	TopicNotExisted          = 1020 //topic不存在
	BundleIdNotExisted       = 1021 //bundle_id不存在
	TopicStartFail           = 1022 //topic开启处理失败
	TopicTypeNotExisted      = 1023 //topic处理类型不存在
	TopicIsNotNull           = 1024 //topic不能为空
	DeviceNotExisted         = 1025 //设备不存在
	StateExisted             = 1027 //状态已存在
	LastMenuNotExisted       = 1028 //上级菜单不存在
	MenuNotExisted           = 1029 //菜单不存在
	UserMenuNotExisted       = 1030 //用户权限不存在
	DeviceIdNotExisted       = 1031 //设备ID不存在
	GoodsDealTypeNotExisted  = 1032 //商品处理类型不存在
	GoodsIdNotExisted        = 1033 //商品不存在
	GoodsBeInDiscount        = 1034 //商品正在打折
	GoodsPayTypeNotExisted   = 1035 //商品可支付类型不存在
	GoodsIdExisted           = 1036 //商品已存在
	OrderIdNotExisted        = 1043 //订单不存在
	GoodsBeNotInDiscount     = 1044 //商品未打折
	NotifyIsNotMatch         = 1045 //会话不匹配
	GoodsIsDiscountRecovery  = 1046 //商品已恢复原价
	InvitationUserNotExisted = 1047 //邀请用户不存在
	//	InvitationUserLevelIsFull = 1048 //邀请用户级数已满
	UserNotAuthorize      = 1049 //用户未授权
	ApplicantIsExisted    = 1050 //申请人已存在
	ApplicantNotExisted   = 1051 //申请人不存在
	CreditOrderNotVaild   = 1052 //订单无效
	RepeatWxWithdraw      = 1053 //微信零钱重复提现
	WxWithdrawAmountError = 1054 //提现金额错误
	WxWithdrawError       = 1055 //微信提现失败
	RepeatSubmission      = 1056 //重复提交
	BundleExisted         = 1057 //bundle已存在
	AuthExisted           = 1058 //权限已存在
	AuthNotExisted        = 1059 //权限不存在
	RoomTypeNotExisted    = 1060 //房型不存在
	RoomTypeExisted       = 1061 //房型已存在
	RoomNoNotExisted      = 1062 //房间不存在
	RoomNoExisted         = 1063 //房间已存在
	RateTypeExisted       = 1064 //房价代码或房价名称已存在
	RateTypeNotExisted    = 1065 //房价代码不存在
	FileNotExisted        = 1066 //文件不存在
	RoomNoInvaild         = 1067 //房间未启用
	ClassExisted          = 1068 //班次已存在
	ClassNotExisted       = 1069 //班次不存在
	CheckTimeError        = 1070 //系统时间与营业时间不匹配
	CurrentClassIsShift   = 1071 //当前班次已交班
	PayPriceError         = 1072 //支付金额错误
	StockNotEnough        = 1073 //存量不足
	DBSaveError           = 1074 //数据存储错误
	DBAddError            = 1075 //数据添加错误
	DBUpdateError         = 1076 //数据更新错误
	DBDeleteError         = 1077 //数据删除错误
	TimeError             = 1078 //时间错误
	OrderInfoError        = 1079 //预定信息错误
	NotVaildError         = 1080 //不允许
	Overdue               = 1081 //已过期
	MaxOverError          = 1082 //超过最大值
	MinOverError          = 1083 //低于最小值
	ExistedError          = 1084 //已存在
	NotBindError          = 1085 //未绑定
	BindError             = 1086 //绑定失败
	CalError              = 1087 //计算错误
)

//消息翻译
var MessageMap = map[int]string{
	NormalMessageId:   "",
	ServerMaintenance: "服务器维护中 请稍后再试",
	AccountDisabled:   "帐号被禁用",
	AppidOverdue:      "appid过期",

	UnknownError:  "未知错误",
	TokenFailure:  "Token失效",
	HTMLSuccess:   "成功",
	BlockingAcess: "禁止访问",

	NewReport: "新消息返回",
	NewHeart:  "心跳消息",

	ParameterInvalid:         "参数无效",
	AppidParameterInvalid:    "授权参数无效",
	EncryptCheckError:        "密文校验失败",
	UserNameDoNotExist:       "用户名不存在或密码错误",
	DuplicateKeyError:        "键值对重复",
	NotFindError:             "未找到",
	InValidOp:                "无效操作",
	InValidAuthorize:         "授权码错误或未找到",
	HasusedError:             "已被使用",
	HasActvError:             "已被激活",
	ActvFailure:              "激活码被禁止使用",
	UserExisted:              "用户已存在",
	VerifyTimeError:          "验证码请求过于平凡",
	MailSendFaild:            "邮箱发送失败",
	SMSSendFaild:             "手机发送失败",
	PhoneParameterError:      "手机号格式有问题",
	VerifyError:              "验证码错误",
	UserNotExisted:           "用户不存在",
	TopicExisted:             "Topic已存在",
	TopicNotExisted:          "Topic不存在",
	BundleIdNotExisted:       "Bundle_id不存在",
	TopicStartFail:           "Topic开启处理失败",
	TopicTypeNotExisted:      "Topic处理类型不存在",
	TopicIsNotNull:           "Topic不能为空",
	DeviceNotExisted:         "设备不存在",
	StateExisted:             "状态已存在",
	LastMenuNotExisted:       "父菜单或根菜单不存在",
	MenuNotExisted:           "菜单不存在",
	UserMenuNotExisted:       "用户权限不存在",
	DeviceIdNotExisted:       "设备ID不存在",
	GoodsDealTypeNotExisted:  "商品处理类型不存在",
	GoodsIdNotExisted:        "商品不存在",
	GoodsBeInDiscount:        "商品正在折扣期",
	GoodsPayTypeNotExisted:   "商品可支付类型不存在",
	GoodsIdExisted:           "商品已存在",
	OrderIdNotExisted:        "订单不存在",
	GoodsBeNotInDiscount:     "商品未打折",
	NotifyIsNotMatch:         "会话不匹配",
	GoodsIsDiscountRecovery:  "商品已恢复原价",
	InvitationUserNotExisted: "邀请用户不存在",
	//	InvitationUserLevelIsFull: "邀请用户级数已满",
	UserNotAuthorize:      "用户未授权",
	ApplicantIsExisted:    "申请人已存在",
	ApplicantNotExisted:   "申请人不存在",
	CreditOrderNotVaild:   "订单无效",
	RepeatWxWithdraw:      "微信零钱一天内多次提现",
	WxWithdrawAmountError: "微信零钱提现金额错误",
	WxWithdrawError:       "微信零钱提现失败",
	RepeatSubmission:      "重复提交",
	BundleExisted:         "bundle已存在",
	AuthExisted:           "权限已存在",
	AuthNotExisted:        "权限不存在",
	RoomTypeNotExisted:    "房型不存在",
	RoomTypeExisted:       "房型已存在",
	RoomNoNotExisted:      "房间不存在",
	RoomNoExisted:         "房间已存在",
	RateTypeExisted:       "房价代码或房价名称已存在",
	RateTypeNotExisted:    "房价代码不存在",
	FileNotExisted:        "文件不存在",
	RoomNoInvaild:         "房间未启用",
	ClassExisted:          "班次已存在",
	ClassNotExisted:       "班次不存在",
	CheckTimeError:        "系统时间与营业时间不匹配",
	CurrentClassIsShift:   "当前班次已完成交班",
	PayPriceError:         "支付金额错误",
	StockNotEnough:        "存量不足",
	DBSaveError:           "数据存储错误",
	DBAddError:            "数据添加错误",
	DBUpdateError:         "数据更新错误",
	DBDeleteError:         "数据删除错误",
	TimeError:             "时间错误",
	OrderInfoError:        "预定信息有误",
	NotVaildError:         "不允许",
	Overdue:               "已过期",
	MaxOverError:          "超过最大值",
	MinOverError:          "低于最小值",
	ExistedError:          "已存在",
	NotBindError:          "未绑定",
	BindError:             "绑定失败",
	CalError:              "计算错误",
}

//MessageBody 消息头
type MessageBody struct {
	State bool        `json:"state"`
	Code  int         `json:"code,omitempty"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

//GetErrorMsg 获取错误消息 参数(int,string)
func GetErrorMsg(errorCode ...interface{}) (msg MessageBody) {
	if len(errorCode) == 0 {
		log.Println("未知")
		msg.State = false
		msg.Code = -1
		return
	}
	msg.State = false
	for _, e := range errorCode {
		switch v := e.(type) {
		case int:
			msg.Code = int(int(v))
			msg.Error = MessageMap[msg.Code]
		case string:
			msg.Error = string(v)
		case error:
			msg.Error = v.Error()
		case interface{}:
			{
				if dev.OnIsDev() {
					msg.Error = onCheckParam(v)
				}
			}
		}
	}

	return
}

func onCheckParam(op interface{}) string {
	//过滤可不填项
	b, _ := json.Marshal(op)

	js, _ := simplejson.NewJson(b)

	mp, _ := js.Map()
	for k, v := range mp {
		tmp := tools.AsString(v)
		if len(tmp) > 0 && tmp != "0" { //过滤
			delete(mp, k)
		}
	}
	//----------------------end

	b, _ = json.Marshal(mp)
	return string(b)
}

//成功消息
func GetSuccessMsg(code ...int) (msg MessageBody) {
	msg.State = true
	if len(code) == 0 {
		msg.Code = NormalMessageId
	} else {
		msg.Code = code[0]
	}

	msg.Error = MessageMap[msg.Code]
	return
}
