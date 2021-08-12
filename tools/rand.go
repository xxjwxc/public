package tools

import (
	"math/rand"
	"time"
)

//生成随机字符串
var _bytes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var r *rand.Rand

// GetRandomString 生成随机字符串
func GetRandomString(n int) string {
	result := []byte{}
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	for i := 0; i < n; i++ {
		result = append(result, _bytes[r.Intn(len(_bytes))])
	}
	return string(result)
}

// GetRangeNumString 生成随机数字字符串
func GetRangeNumString(n int) string {
	var _bytes = []byte("0123456789")
	var r *rand.Rand

	result := []byte{}
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	for i := 0; i < n; i++ {
		result = append(result, _bytes[r.Intn(len(_bytes))])
	}
	return string(result)
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
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// GetRandInt 生成随机整数 digit：位数
func GetRandInt(min int, max int) int {
	if min > max {
		min = 0
		max = 0
	}
	if max == min {
		return min
	}
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
