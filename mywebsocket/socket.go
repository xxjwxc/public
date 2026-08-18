package mywebsocket

import (
	"net/http"
	"sync"
	"time"

	gws "github.com/gorilla/websocket"
	"github.com/xxjwxc/public/myglobal"
	"github.com/xxjwxc/public/mylog"
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

	// heartbeat
	alive  bool
	pongMu sync.Mutex
	hbDone chan struct{}
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
	s.StopHeartbeat()
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.conn.Close()
}

// Ping 发送 WebSocket Ping 控制帧
func (s *Socket) Ping() error {
	return s.conn.WriteControl(gws.PingMessage, []byte{}, time.Now().Add(5*time.Second))
}

// StartHeartbeat 启动心跳探活机制
// interval: Ping 发送间隔, 超时未收到 Pong 则判定连接断开并关闭连接
func (s *Socket) StartHeartbeat(interval, timeout time.Duration) {
	s.StopHeartbeat()

	s.pongMu.Lock()
	s.alive = true
	s.pongMu.Unlock()

	// 设置 Pong 处理器（在 ReadMessage 调用时自动触发）
	s.conn.SetPongHandler(func(appData string) error {
		s.pongMu.Lock()
		s.alive = true
		s.pongMu.Unlock()
		return nil
	})

	s.hbDone = make(chan struct{})
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.pongMu.Lock()
				if !s.alive {
					s.pongMu.Unlock()
					mylog.Infof("[websocket] heartbeat timeout, closing connection [%v]", s.SessionId)
					s.conn.Close()
					return
				}
				s.alive = false
				s.pongMu.Unlock()

				if err := s.Ping(); err != nil {
					mylog.Infof("[websocket] ping failed, closing connection [%v]: %v", s.SessionId, err)
					return
				}
			case <-s.hbDone:
				return
			}
		}
	}()
}

// StopHeartbeat 停止心跳探活
func (s *Socket) StopHeartbeat() {
	if s.hbDone != nil {
		close(s.hbDone)
		s.hbDone = nil
	}
}
