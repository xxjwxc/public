package myssh

import (
	"fmt"
)

// MyWriter io.Writer
type MyWriter struct {
	channel chan string
}

// NewWriter new writer
func NewWriter() *MyWriter {
	w := &MyWriter{
		channel: make(chan string),
	}
	return w
}

func (w *MyWriter) Write(p []byte) (n int, err error) {
	w.channel <- string(p)
	return len(p), err
}

// Consume 消费
func (w *MyWriter) Consume() string {
	return <-w.channel
}

// Close 关闭
func (w *MyWriter) Close() {
	close(w.channel)
}

// Run 消费
func (w *MyWriter) Run() {
	go func() {
		for res := range w.channel {
			fmt.Print(res)
			// 以下去掉命令行显示
			// index := strings.Index(res, "\n")
			// fmt.Print(res[index+1:])
		}
	}()
}
