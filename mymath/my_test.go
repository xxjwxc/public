package mymath

import (
	"fmt"
	"testing"
	"time"
)

func Test_order(t *testing.T) {
	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("aaaa")
		ch <- 123
	}()

	h := <-ch
	fmt.Println("bbbb")
	fmt.Println(h)

	fmt.Println(Gcd(9, 21))
	fmt.Println(17 * 19)
	fmt.Println(Lcm(17, 19))
}
