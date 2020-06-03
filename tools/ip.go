package tools

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/xxjwxc/public/mylog"
)

// GetWwwIP 获取公网IP地址
func GetWwwIP() (exip string) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(b))
}

// GetLocalIP 获取内网ip
func GetLocalIP() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		mylog.Error(err)
		return
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return
}

// GetClientIP 获取用户ip
func GetClientIP(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return
}
