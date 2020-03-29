package myi18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/xxjwxc/public/tools"
	"golang.org/x/text/language"
)

func init() {
	ReSet()
}

// SetLocalLG 设置本地语言包 （空 将获取系统值）
func SetLocalLG(lg string) {
	if len(lg) == 0 {
		lg = tools.GetLocalSystemLang(true)
	}
	// return the new localizer that can be used to translate text
	tr = i18n.NewLocalizer(i18nBundle, lg)
}

// AddMessages 添加语言
func AddMessages(tag language.Tag, messages ...*i18n.Message) error {
	return i18nBundle.AddMessages(tag, messages...)
}

// AddKV 添加语言
func AddKV(tag language.Tag, k, v string) error {
	return i18nBundle.AddMessages(tag, &i18n.Message{
		ID:    k,
		Other: v,
	})
}

// ReSet 重置
func ReSet() {
	i18nBundle = i18n.NewBundle(language.English)
	SetLocalLG("") // default
}

// Get 获取值
func Get(Key string) string {
	if tr != nil {
		return tr.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: Key,
			},
		})
	}

	return Key
}
