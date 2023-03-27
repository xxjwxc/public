package mywebsocket

/*
	说明：第一个包 初始化client唯一id。消息id为100
*/
import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xxjwxc/public/message"
)

type ClientBody struct {
	Code int         `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type WSMessageType uint8

const (
	WS_Login WSMessageType = 1 // 链接
	WS_Close WSMessageType = 2 // 断线
)

type HandlerReadFunc func(messageType int, p []byte, err error) //ID , MESSAGEID,内容

type MyWebSocket struct {
	conn *websocket.Conn
	done chan struct{}
}

//websocket 初始化接口
/*
addr 地址,
path 域
handlerFunc 读到的消息回调
stateFunc 连接消息回调
timeOut 读取超时回调(0则永不超时)
*/
func NewWebSocket(url string, handlerFunc HandlerReadFunc, requestHeader http.Header, timeOut time.Duration) (*MyWebSocket, error) {
	myWebSocket := &MyWebSocket{
		conn: nil,
		done: make(chan struct{}),
	}

	// 	key := `-----BEGIN CERTIFICATE-----
	// MIIBkTCCAT+gAwIBAgIQC+AfyVkg13ejxOIycKc6kzAKBggqhkjOPQQDAjASMRAw
	// DgYDVQQKEwdBY21lIENvMB4XDTIzMDMxNTEzMTEzNVoXDTI0MDMxNDEzMTEzNVow
	// EjEQMA4GA1UEChMHQWNtZSBDbzBOMBAGByqGSM49AgEGBSuBBAAhAzoABO5vPodF
	// 2Gtpxm3e7uXQGbiA3d+hHR0KydxTxqZwnS5lAKO/EHYwHbQYrgI8jDuKi/ZRH3HN
	// l+AXo4GBMH8wDgYDVR0PAQH/BAQDAgKkMBMGA1UdJQQMMAoGCCsGAQUFBwMBMA8G
	// A1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFMLLs994lljKRKNUvsNTrmCsKvQVMCgG
	// A1UdEQQhMB+CHXdlYmNhc3QzLXdzLXdlYi1sZi5kb3V5aW4uY29tMAoGCCqGSM49
	// BAMCA0AAMD0CHQDzAgEPWy0AS2ovCju0r8IOII9rtiLqagHKL7ykAhwryu3E7hU1
	// Jhhe7V5DyPrX0aPBAHYIP6NKHHcx
	// -----END CERTIFICATE-----
	// 	`
	// 	certs := x509.NewCertPool()
	// 	ok := certs.AppendCertsFromPEM([]byte(key))
	// 	if !ok {
	// 		log.Fatal("failed to parse root certificate")
	// 	}
	var dialer = &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  timeOut,
		EnableCompression: true,
		// TLSClientConfig:   &tls.Config{RootCAs: nil, InsecureSkipVerify: true},
	}
	// dialer.TLSClientConfig = &tls.Config{RootCAs: certs, InsecureSkipVerify: true}
	var err error
	var resp *http.Response
	myWebSocket.conn, resp, err = dialer.Dial(url, requestHeader)
	if err != nil {
		fmt.Println(resp)
		return nil, err
	}

	go func() {
		defer close(myWebSocket.done)
		for {
			messageType, message, err := myWebSocket.conn.ReadMessage()
			handlerFunc(messageType, message, err)
		}
	}()

	return myWebSocket, nil
}

func (wss *MyWebSocket) SendMessage(p []byte) error {
	if wss.conn != nil {
		return wss.conn.WriteMessage(websocket.TextMessage, p)
	}

	return message.GetError(message.EmptyError)
}

func (wss *MyWebSocket) Close() {
	wss.done <- struct{}{}
	if wss.conn != nil {
		wss.conn.Close()
	}
}
