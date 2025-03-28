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
	}, nil
}

type Socket struct {
	SessionId string
	conn      *gws.Conn
	sync.Mutex
}

func (s *Socket) ID() string {
	return s.SessionId
}

func (s *Socket) WriteMessage(byteMessage []byte) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.conn.WriteMessage(gws.BinaryMessage, byteMessage)
}

func (s *Socket) ReadMessage() (messageType int, p []byte, err error) {
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
