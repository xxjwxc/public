package mywebsocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gws "github.com/gorilla/websocket"
)

// setupTestServer 创建测试 WebSocket 服务端
func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("upgrade error: %v", err)
			return
		}
		defer conn.Close()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))
}

// dialClient 连接测试服务端
func dialClient(t *testing.T, server *httptest.Server) *gws.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := gws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	return conn
}

func TestPing(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	conn := dialClient(t, server)
	defer conn.Close()

	// 设置 pong handler
	pongCh := make(chan string, 1)
	conn.SetPongHandler(func(appData string) error {
		pongCh <- appData
		return nil
	})

	sock := &Socket{conn: conn, d: time.Second * 30}
	if err := sock.Ping(); err != nil {
		t.Fatalf("Ping() error: %v", err)
	}

	// 需要 ReadMessage 来触发 pong 帧处理
	go func() {
		conn.ReadMessage()
	}()

	select {
	case <-pongCh:
		t.Log("received pong")
	case <-time.After(3 * time.Second):
		t.Fatal("pong not received within timeout")
	}
}

func TestStartHeartbeat(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	conn := dialClient(t, server)
	defer conn.Close()

	sock := &Socket{
		SessionId: "test-hb",
		conn:      conn,
		d:         time.Second * 30,
	}

	sock.StartHeartbeat(500*time.Millisecond, 2*time.Second)
	defer sock.StopHeartbeat()

	// 后台持续 ReadMessage 以触发 pong handler
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)

	err := sock.Ping()
	if err != nil {
		t.Fatalf("connection should be alive after heartbeat, got: %v", err)
	}
	t.Log("heartbeat kept connection alive")
}

func TestStopHeartbeat(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	conn := dialClient(t, server)
	defer conn.Close()

	sock := &Socket{
		SessionId: "test-stop",
		conn:      conn,
		d:         time.Second * 30,
	}

	sock.StartHeartbeat(200*time.Millisecond, 1*time.Second)
	sock.StopHeartbeat()

	if sock.hbDone != nil {
		t.Fatal("hbDone should be nil after StopHeartbeat")
	}

	// 连续调用不应 panic
	sock.StopHeartbeat()
	t.Log("StopHeartbeat is idempotent")
}

func TestHeartbeatTimeout(t *testing.T) {
	// 服务端不读取消息，使 pong 无法被处理
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		// 故意不 ReadMessage
		time.Sleep(10 * time.Second)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := gws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}

	sock := &Socket{
		SessionId: "test-timeout",
		conn:      conn,
		d:         time.Second * 30,
	}

	// 心跳: 200ms ping, 超时后 conn.Close() 会被调用
	sock.StartHeartbeat(200*time.Millisecond, 500*time.Millisecond)
	defer sock.StopHeartbeat()

	// 等待超时触发
	time.Sleep(2 * time.Second)

	// 验证连接已被关闭：尝试写入应该失败
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	err = conn.WriteMessage(gws.TextMessage, []byte("test"))
	if err == nil {
		t.Fatal("connection should have been closed due to heartbeat timeout")
	}
	t.Logf("heartbeat timeout correctly closed connection: %v", err)
}
