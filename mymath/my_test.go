package mymath

import "testing"

func TestGcd(t *testing.T) {
	cases := []struct {
		x, y, want int64
	}{
		{9, 21, 3},
		{17, 19, 1},
		{12, 12, 12},
		{0, 5, 5},   // 原实现会 panic: integer divide by zero
		{5, 0, 5},   // 原实现会 panic: integer divide by zero
		{0, 0, 0},   // 原实现会 panic: integer divide by zero
		{-9, 21, 3}, // 负数取绝对值后再求
		{-9, -21, 3},
	}
	for _, c := range cases {
		if got := Gcd(c.x, c.y); got != c.want {
			t.Errorf("Gcd(%d,%d) = %d, 期望 %d", c.x, c.y, got, c.want)
		}
	}
}

func TestLcm(t *testing.T) {
	cases := []struct {
		x, y, want int64
	}{
		{4, 6, 12},
		{17, 19, 323},
		{0, 5, 0},
		{5, 0, 0},
	}
	for _, c := range cases {
		if got := Lcm(c.x, c.y); got != c.want {
			t.Errorf("Lcm(%d,%d) = %d, 期望 %d", c.x, c.y, got, c.want)
		}
	}
}
