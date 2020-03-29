package fractional

import (
	"fmt"
	"testing"
)

func Test_fal(t *testing.T) {
	tmp := Model(7, 12)
	tmp1 := Model(1, 12)
	fmt.Println(tmp.Add(tmp1))

	tmp = Model(1, 4)
	tmp1 = Model(1, 3)
	fmt.Println(tmp.Sub(tmp1))

	tmp = Model(3, 4)
	tmp1 = Model(2, 3)
	fmt.Println(tmp.Mul(tmp1))

	tmp = Model(3, 4)
	tmp1 = Model(2, 3)
	fmt.Println(tmp.Div(tmp1))

	tmp = Model(1, 3)
	fmt.Println(tmp.Verdict())

}
