package mykdniao

import (
	"fmt"
	"testing"

	"github.com/xxjwxc/public/tools"
)

func Test_kdn(t *testing.T) {
	kdn := New("1111111", "11111111-1111-1111-1111-11111111111111")
	result := kdn.GetLogisticsTrack("4304678557725", "YD")
	fmt.Printf(tools.JSONDecode(result))
}
