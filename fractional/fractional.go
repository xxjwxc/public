/*
Package fractional 分数运算相关
*/
package fractional

import (
	"fmt"

	"github.com/xxjwxc/public/mymath"
)

//FAL 分数
type FAL struct {
	Nume int64 //numerator 分子
	Deno int64 //denominator 分母 (一定不为0)
}

//Model Create a score (molecular, denominator) with a denominator default of 1 创建一个分数(分子，分母)，分母默认为1:
func Model(nd ...int64) *FAL {
	var f FAL
	if len(nd) == 1 {
		f.Nume = nd[0]
		f.Deno = 1
	} else if len(nd) == 2 {
		f.Nume = nd[0]
		f.Deno = nd[1]
	}

	if f.Deno == 0 { //分母为0
		panic(fmt.Sprintf("fractional init error. if denominator can't zero."))
	}

	return &f
}

//阔张
func (s *FAL) broad(lcm int64) {
	s.Nume = s.Nume * (lcm / s.Deno)
	s.Deno = lcm
}

//压缩 整理
func (s *FAL) offset() {
	lcm := mymath.Gcd(s.Nume, s.Deno)

	s.Nume /= lcm
	s.Deno /= lcm
}

//Add 分数加法
func (s *FAL) Add(f *FAL) *FAL {
	//获取最小公倍数
	lcm := mymath.Lcm(f.Deno, s.Deno)
	s.broad(lcm)
	f.broad(lcm)

	s.Nume += f.Nume
	s.offset()
	return s
}

//Sub 分数减法
func (s *FAL) Sub(f *FAL) *FAL {
	//获取最小公倍数
	lcm := mymath.Lcm(s.Deno, f.Deno)
	s.broad(lcm)
	f.broad(lcm)

	s.Nume -= f.Nume
	s.offset()
	return s
}

//Mul 乘法
func (s *FAL) Mul(f *FAL) *FAL {
	s.Deno *= f.Deno
	s.Nume *= f.Nume
	s.offset()
	return s
}

//Div 乘法
func (s *FAL) Div(f *FAL) *FAL {
	tmp := Model(f.Deno, f.Nume)
	s.Mul(tmp)
	s.offset()
	return s
}

//Verdict 计算结果
func (s *FAL) Verdict() float64 {
	return float64(s.Nume) / float64(s.Deno)
}
