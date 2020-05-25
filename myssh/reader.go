package myssh

import (
	"bufio"
	"fmt"
	"os"
	"sync/atomic"
)

// MyReader io reader
type MyReader struct {
	channel chan string
	isClose int32
}

// NewReader new io reader
func NewReader() *MyReader {
	r := &MyReader{
		channel: make(chan string),
		isClose: 0,
	}
	return r
}

func (r *MyReader) Read(p []byte) (n int, err error) {
	fmt.Println("into Read...")
	fmt.Println("=======:", string(p))
	cmd := <-r.channel
	tmp := []byte(cmd + "\n")
	for i, v := range tmp {
		p[i] = v
	}
	fmt.Println("leave Read.")
	return len(tmp), err
}

// Push push one string
func (r *MyReader) Push(src string) {
	r.channel <- src
}

// ListenStdin 监听cmd 输入
func (r *MyReader) ListenStdin() {
	go func() {
		f := bufio.NewReader(os.Stdin) //读取输入的内容
		for {
			var str, arg string
			Input, _ := f.ReadString('\n') //定义一行输入的内容分隔符。
			fmt.Sscan(Input, &str, &arg)   //将Input
			if len(arg) > 0 {
				arg = " " + arg
			}

			fmt.Println("yes:::::", str+arg)
			if atomic.LoadInt32(&r.isClose) == 1 {
				fmt.Println("outtttttt")
				break
			} else {
				r.channel <- str + arg
			}
		}
	}()
}

// Close 关闭
func (r *MyReader) Close() {
	atomic.StoreInt32(&r.isClose, 1)
	close(r.channel)
}
