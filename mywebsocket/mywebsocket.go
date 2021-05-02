package mywebsocket

/*
	说明：第一个包 初始化client唯一id。消息id为100
*/
import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/xxjwxc/public/mycache"
	"github.com/xxjwxc/public/mylog"

	"golang.org/x/net/websocket"
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

type HandlerReadFunc func(string, string, ClientBody, *websocket.Conn) //ID , MESSAGEID,内容
type HandlerStateFunc func(string, string, WSMessageType)              //状态发声改变回调

var mutex sync.Mutex

/*
	写数据
	pathExp:请求的根路径
	clientid:发送人id
	body:发送内容
*/
func WriteData(pathExp string, clientid string, body ClientBody) bool {

	wb, err := json.Marshal(body)
	if err != nil {
		mylog.Debug("error:", err.Error())
		return false
	}

	cache := mycache.NewCache("websocket" + pathExp)
	var tp []*websocket.Conn
	e := cache.Value(clientid, &tp)

	if e == nil {
		b_r := false
		tmp := tp
		for i := 0; i < len(tmp); i++ {
			if _, err = tmp[i].Write(wb); err != nil {
				mylog.Debug("Can't send", err.Error())
				defer tmp[i].Close()
			} else {
				b_r = true
			}
		}
		return b_r
	} else {
		mylog.Debug("not find client:" + clientid)
		return false
	}
}

//websocket 初始化接口
/*
pathExp 域,
handlerFunc 读到的消息回调
stateFunc 连接消息回调
timeOut 读取超时回调(0则永不超时)
isMult 是否多用户登录
*/
func InitWebSocket(pathExp string, handlerFunc HandlerReadFunc, stateFunc HandlerStateFunc, timeOut time.Duration, isMult bool) *rest.Route {
	wsHandler := websocket.Handler(func(ws *websocket.Conn) {
		mylog.Debug("enter")
		var err error
		var msg = make([]byte, 1024)
		var n int
		n, err = ws.Read(msg)

		if err != nil {
			mylog.Debug("ws:close")
			return
		}

		var clientid string
		var clientBody ClientBody

		err = json.Unmarshal(msg[:n], &clientBody)
		if err != nil {
			mylog.Debug("Unmarshal:" + err.Error())
			return
		} else {
			if clientBody.Code != 100 {
				mylog.Debug("messageid error")
				return
			}
			clientid = clientBody.Data.(string)
			//保存缓存
			cache := mycache.NewCache("websocket" + pathExp)
			var tmp []*websocket.Conn

			mutex.Lock()

			var tp []*websocket.Conn
			b := cache.Value(clientid, &tp)
			if b == nil && isMult { //多用户
				tmp = tp
			}

			tmp = append(tmp, ws)
			cache.Add(clientid, tmp, 2*time.Hour) //2小时过期
			mutex.Unlock()

			//------------------end
			mylog.Debug("init success:" + clientid)
			if stateFunc != nil {
				stateFunc(pathExp, clientid, WS_Login)
			}
		}

		ch := make(chan bool, 1)
		if timeOut > 0 {
			go func(ws *websocket.Conn) {
				var after <-chan time.Time
			loop:
				after = time.After(timeOut)
				for {
					select {
					case b := <-ch: //继续下一个等待
						if !b {
							break
						} else {
							goto loop
						}
					case <-after: //超时处理
						mylog.Info("time out:" + clientid)
						ws.Close()
						break
					}
				}
			}(ws)
		}

		for {
			n, err = ws.Read(msg)
			if err != nil {
				if timeOut > 0 {
					ch <- false
				}
				mylog.Debug("ws:close")
				break
			} else if timeOut > 0 {
				ch <- true
			}
			var body ClientBody
			err = json.Unmarshal(msg[:n], &body)
			if err != nil {
				mylog.Debug("r:" + err.Error())
			} else {
				if handlerFunc != nil {
					handlerFunc(pathExp, clientid, body, ws)
				}
			}
		}

		//删除缓存
		cache := mycache.NewCache("websocket" + pathExp)
		var tmp []*websocket.Conn

		mutex.Lock()

		var tp []*websocket.Conn
		b := cache.Value(clientid, &tp)
		if b == nil {
			tmp = tp
		}
		i := 0
		for ; i < len(tmp); i++ {
			if tmp[i] == ws {
				tmp = append(tmp[:i], tmp[i+1:]...) // 最后面的“...”不能省略
				break
			}
		}
		if i == 0 || len(tmp) == 0 || !isMult {
			cache.Delete(clientid)
			mylog.Debug("delete all: " + clientid)
		} else {
			cache.Add(clientid, tmp, 2*time.Hour) //2小时过期
			mylog.Debug("delete one: " + clientid)
		}

		mutex.Unlock()
		if stateFunc != nil {
			stateFunc(pathExp, clientid, WS_Close)
		}
		//------------------end

	})

	return rest.Get(pathExp, func(w rest.ResponseWriter, r *rest.Request) {
		//mylog.Debug("-------------")
		wsHandler.ServeHTTP(w.(http.ResponseWriter), r.Request)
	})
}
