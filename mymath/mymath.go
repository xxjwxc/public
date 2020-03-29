package mymath

import "math"

//Gcd 最大公约数:(辗转相除法)
func Gcd(x, y int64) int64 {
	x = int64(math.Abs(float64(x)))
	y = int64(math.Abs(float64(y)))

	var tmp int64
	for {
		tmp = (x % y)
		if tmp > 0 {
			x = y
			y = tmp
		} else {
			return y
		}
	}
}

//Lcm 最小公倍数:((x*y)/最大公约数)
func Lcm(x, y int64) int64 {
	return (x * y) / Gcd(x, y)
}
