package weixin

import (
	"io"
	"net/http"
	"time"
)

// DefaultTimeout 微信接口请求的默认超时时间。
// 直接使用 http.Get / http.Post 走的是 DefaultClient(无超时),网络异常时请求可能永久挂起。
var DefaultTimeout = 30 * time.Second

// defaultClient 带超时的 HTTP 客户端,底层复用 DefaultTransport 的连接池
func defaultClient() *http.Client {
	return &http.Client{Timeout: DefaultTimeout}
}

// httpGet 带超时的 GET 请求
func httpGet(url string) (*http.Response, error) {
	return defaultClient().Get(url)
}

// httpPost 带超时的 POST 请求
func httpPost(url, contentType string, body io.Reader) (*http.Response, error) {
	return defaultClient().Post(url, contentType, body)
}
