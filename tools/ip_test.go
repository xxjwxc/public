package tools

import (
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	t.Logf("Local IP: %s", ip)
	// 验证IP不为空且不是链路本地地址
	if ip == "" {
		t.Error("GetLocalIP returned empty string")
	}
	// 验证IP不是169.254开头的链路本地地址
	if len(ip) >= 7 && ip[:7] == "169.254" {
		t.Errorf("GetLocalIP returned link-local address: %s", ip)
	}
}