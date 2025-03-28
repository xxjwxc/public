package mywebsocket

import (
	"sync"

	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mylog"
)

var (
	clients sync.Map

	joinChannel = make(chan *Socket)
	quitChannel = make(chan *Socket)
)

func init() {
	go connect()
}
func connect() {
	defer func() {
		if err := recover(); err != nil {
			mylog.Errorf("recover: %v", err)
		}
		mylog.Errorf("websocket connect goroutine exit!!!")
	}()

	mylog.Infof("connect goroutine started ...")
	for {
		select {
		case cli := <-joinChannel:
			mylog.Infof("socket join: %v", cli.ID())
			clients.Store(cli.ID(), cli)
		case cli := <-quitChannel:
			mylog.Infof("socket quit: %v", cli.ID())
			clients.Delete(cli.ID())
		}
	}
}

// 添加socket客户端
func AddSocketClient(cli *Socket) {
	joinChannel <- cli
}

// 删除socket客户端
func DelSocketClient(cli *Socket) {
	quitChannel <- cli
}

// SendOneMessageFromId 发送消息给所有客户端
func SendOneMessageFromId(clientId string, byteMessage []byte) error {
	cli, ok := clients.Load(clientId)
	if !ok {
		return message.GetError(message.NotFindError)
	}
	return cli.(*Socket).WriteMessage(byteMessage)
}

func SendAllMessage(byteMessage []byte) {
	clients.Range(func(key, value interface{}) bool {
		value.(*Socket).WriteMessage(byteMessage)
		return true
	})
}

func Length() int {
	count := 0
	clients.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
