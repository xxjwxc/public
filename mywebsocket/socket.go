package mywebsocket

import (
	"net/http"
	"sync"
	"time"

	gws "github.com/gorilla/websocket"
	"github.com/xxjwxc/public/myglobal"
)

var (
	upgrader = gws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: time.Second * 10,
	}
)
var node *myglobal.NodeInfo

func init() {
	node = myglobal.GetNode()
}

func NewSocketUpgrader(sessionId string, w http.ResponseWriter, r *http.Request) (*Socket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	if sessionId == "" {
		sessionId = node.GetIDStr()
	}

	return &Socket{
		SessionId: sessionId,
		conn:      conn,
		d:         time.Second * 30,
	}, nil
}

type Socket struct {
	sync.Mutex
	SessionId string
	conn      *gws.Conn
	d         time.Duration
}

func (s *Socket) ID() string {
	return s.SessionId
}

func (s *Socket) SetTimeOut(d time.Duration) {
	s.d = d
}

func (s *Socket) WriteMessage(byteMessage []byte) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.conn.WriteMessage(gws.BinaryMessage, byteMessage)
}

func (s *Socket) WriteText(msg string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.conn.WriteMessage(gws.TextMessage, []byte(msg))
}

func (s *Socket) ReadMessage() (messageType int, p []byte, err error) {
	// 设置读取超时（示例：5秒）
	if err = s.conn.SetReadDeadline(time.Now().Add(s.d)); err != nil {
		// 处理错误
		return
	}
	return s.conn.ReadMessage()
}

func (s *Socket) Close() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.conn.Close()
}

// func (s *Socket) Upgrade(w http.ResponseWriter, r *http.Request) (*gws.Conn, error) {
// 	return upgrader.Upgrade(w, r, nil)
// }
