package server

import (
	"fmt"
	"testing"
)

func CallBack() {
	fmt.Println("aa")
}

func TestDomain(t *testing.T) {
	On("n", "dn", "d").Start(CallBack)
}
