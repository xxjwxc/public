package tools

import (
	"math/rand"
)

// 随机字符串使用的字符集
var _bytes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randomString 从给定字符集中随机取 n 个字符
// math/rand 的顶层函数是并发安全的,且 Go 1.20+ 起全局随机源会自动随机播种,
// 因此无需自行维护 *rand.Rand 实例(原实现共享 *rand.Rand 会产生数据竞争)
func randomString(chars []byte, n int) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, n)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// GetRandomString 生成随机字符串
func GetRandomString(n int) string {
	return randomString(_bytes, n)
}

// GetRangeNumString 生成随机数字字符串
func GetRangeNumString(n int) string {
	return randomString([]byte("0123456789"), n)
}

// GetRangeNum 生成随机整数 digit：位数
func GetRangeNum(digit int) int {
	var max, min int = 1, 1
	if digit > 0 {
		for i := 0; i < digit; i++ {
			max = max * 10
		}
		for i := 0; i < digit-1; i++ {
			min = min * 10
		}
	}
	if max == min {
		return min
	}
	// 不再调用 rand.Seed:秒级种子会让同一秒内的调用返回完全相同的值,
	// 且会覆盖全局随机源。顶层函数在 Go 1.20+ 已自动随机播种。
	return rand.Intn(max-min) + min
}

// GetRandInt 生成随机整数 digit：位数
func GetRandInt(min int, max int) int {
	if min > max {
		return 0
	}
	if max == min {
		return min
	}
	return rand.Intn(max-min) + min
}
