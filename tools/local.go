package tools

import (
	"os"
	"strings"
)

// Getenv 获取本地系统变量
func Getenv(key string) string {
	return os.Getenv(key)
}

// GetLocalSystemLang 获取本地语言 (like:zh_CN.UTF-8)(simple:zh)
func GetLocalSystemLang(isSimple bool) (locale string) {
	locale = Getenv("LC_ALL")
	if locale == "" {
		locale = Getenv("LANG")
	}
	if isSimple {
		locale, _ = splitLocale(locale)
	}
	if len(locale) == 0 {
		locale = "zh"
	}
	return
}

func splitLocale(locale string) (string, string) {
	formattedLocale := strings.Split(locale, ".")[0]
	formattedLocale = strings.Replace(formattedLocale, "-", "_", -1)

	pieces := strings.Split(formattedLocale, "_")
	language := pieces[0]
	territory := ""
	if len(pieces) > 1 {
		territory = strings.Split(formattedLocale, "_")[1]
	}
	return language, territory
}
