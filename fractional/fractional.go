/*
Package fractional 分数运算相关
*/
package fractional

import (
	"github.com/xxjwxc/public/mymath"
)

// FAL 分数
type FAL struct {
	Nume int64 //numerator 分子
	Deno int64 //denominator 分母 (一定不为0)
}

// Model Create a score (molecular, denominator) with a denominator default of 1 创建一个分数(分子，分母)，分母默认为1:
func Model(nd ...int64) *FAL {
	var f FAL
	switch len(nd) {
	case 1:
		f.Nume = nd[0]
		f.Deno = 1
	case 2:
		f.Nume = nd[0]
		f.Deno = nd[1]
	default:
		panic("fractional init error. params must be (numerator[, denominator])")
	}

	if f.Deno == 0 { //分母为0
		panic("fractional init error. denominator can't be zero")
	}

	return &f
}

// 阔张
func (s *FAL) broad(lcm int64) {
	s.Nume = s.Nume * (lcm / s.Deno)
	s.Deno = lcm
}

// 压缩 整理
func (s *FAL) offset() {
	lcm := mymath.Gcd(s.Nume, s.Deno)
	if lcm == 0 { // 理论上不会发生,兜底避免除零 panic
		return
	}

	s.Nume /= lcm
	s.Deno /= lcm
}

// Add 分数加法
func (s *FAL) Add(f *FAL) *FAL {
	//获取最小公倍数
	lcm := mymath.Lcm(f.Deno, s.Deno)
	if lcm == 0 {
		panic("fractional add error. denominator can't be zero")
	}
	// 用副本参与通分,避免修改入参 f(原实现会污染调用方数据)
	nume := f.Nume * (lcm / f.Deno)
	s.broad(lcm)

	s.Nume += nume
	s.offset()
	return s
}

// Sub 分数减法
func (s *FAL) Sub(f *FAL) *FAL {
	//获取最小公倍数
	lcm := mymath.Lcm(s.Deno, f.Deno)
	if lcm == 0 {
		panic("fractional sub error. denominator can't be zero")
	}
	// 用副本参与通分,避免修改入参 f(原实现会污染调用方数据)
	nume := f.Nume * (lcm / f.Deno)
	s.broad(lcm)

	s.Nume -= nume
	s.offset()
	return s
}

// Mul 乘法
func (s *FAL) Mul(f *FAL) *FAL {
	s.Deno *= f.Deno
	s.Nume *= f.Nume
	s.offset()
	return s
}

// Div 乘法
func (s *FAL) Div(f *FAL) *FAL {
	tmp := Model(f.Deno, f.Nume)
	s.Mul(tmp)
	s.offset()
	return s
}

// Verdict 计算结果
func (s *FAL) Verdict() float64 {
	return float64(s.Nume) / float64(s.Deno)
}
