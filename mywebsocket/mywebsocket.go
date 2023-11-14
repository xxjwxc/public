package mywebsocket

/*
	说明：第一个包 初始化client唯一id。消息id为100
*/
import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mylog"
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

	var err error
	var resp *http.Response
	myWebSocket.conn, resp, err = websocket.DefaultDialer.Dial(url, requestHeader)
	if err != nil {
		mylog.Error(resp)
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

func (wss *MyWebSocket) SendMessage(messageType int, p []byte) error {
	if wss.conn != nil {
		return wss.conn.WriteMessage(messageType, p)
	}

	return message.GetError(message.EmptyError)
}

func (wss *MyWebSocket) Close() {
	wss.done <- struct{}{}
	if wss.conn != nil {
		wss.conn.Close()
	}
}
