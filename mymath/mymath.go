package mymath

// abs64 取绝对值(整数运算,避免 float64 转换在超过 2^53 时丢失精度)
func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// Gcd 最大公约数:(辗转相除法)
// 当任一参数为 0 时返回另一个参数的绝对值,两者均为 0 时返回 0(不会发生除零 panic)
func Gcd(x, y int64) int64 {
	x = abs64(x)
	y = abs64(y)
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

// Lcm 最小公倍数:((x*y)/最大公约数)
// 任一参数为 0 时返回 0
func Lcm(x, y int64) int64 {
	if x == 0 || y == 0 {
		return 0
	}
	return (x / Gcd(x, y)) * y // 先除后乘,降低溢出概率
}
